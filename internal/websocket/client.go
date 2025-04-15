package websocket

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qianmianyao/parchment-server/internal/websocket/message_type"
	"github.com/qianmianyao/parchment-server/pkg/encryption"
	"github.com/qianmianyao/parchment-server/pkg/global"
)

// WebSocket 连接相关的常量定义。
const (
	// writeWait 是允许向对端写入消息的最大等待时间。
	writeWait = 10 * time.Second

	// pongWait 是允许从对端读取下一个 Pong 消息的最大等待时间。
	pongWait = 60 * time.Second

	// pingPeriod 是向对端发送 Ping 消息的时间间隔, 必须小于 pongWait。
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize 是允许从对端接收的 WebSocket 消息的最大大小（字节）。
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// upgrader 用于将 HTTP 连接升级为 WebSocket 连接。
// 注意：生产环境中 CheckOrigin 应实现更严格的安全策略。
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 生产环境需要更严格的来源检查
		return true
	},
}

// Client 代表一个 WebSocket 客户端。
type Client struct {
	hub      *Hub            // 所属的 Hub。
	conn     *websocket.Conn // WebSocket 连接。
	send     chan []byte     // 发送消息的缓冲通道。
	uuid     string          // 客户端唯一标识符 (User ID)。
	username string          // 客户端用户名。
}

// readPump 从 WebSocket 连接读取消息并传递给 Hub 处理。
// 同时处理连接关闭和 Pong 消息以维持连接。
func (c *Client) readPump() {
	// 确保在退出时注销客户端并关闭连接。
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Error closing connection for %s: %v", c.uuid, err))
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 设置 Pong 消息处理器，收到 Pong 时更新读取截止时间。
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Warn(fmt.Sprintf("Unexpected websocket close for %s: %v", c.uuid, err))
			}
			break // 退出循环
		}
		// 清理消息：移除首尾空格，换行符替换为空格。
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// global.Logger.Debug(fmt.Sprintf("Received message from %s: %s", c.uuid, string(message))) // 可选调试日志

		// 解析消息。
		_, envelope, err := message_type.ParseMessage(message)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Failed to parse message from %s: %v", c.uuid, err))
			continue
		}

		// 根据消息目标路由。
		if envelope.Destination != "all" && envelope.Destination != "" {
			// 发送给特定客户端或房间。
			c.hub.SendToSpecificClient(envelope.Source.Uid, envelope.Destination, message)
		} else {
			// 广播消息（当前注释掉）。
			// c.hub.broadcast <- message
		}
	}
}

// writePump 将 `send` 通道中的消息写入 WebSocket 连接。
// 同时通过定期发送 Ping 消息维持连接活跃。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	// 确保在退出时停止定时器并关闭连接。
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Error closing connection for %s: %v", c.uuid, err))
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// send 通道已关闭，通知对端关闭。
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 获取写入器，同一时间只允许一个写入器活跃。
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to get next writer for %s: %v", c.uuid, err))
				return
			}

			// 写入当前消息。
			_, err = w.Write(message)
			if err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to write message for %s: %v", c.uuid, err))
				_ = w.Close() // 即使写入失败，也尝试关闭写入器
				return
			}

			// 检查并写入 `send` 通道中的排队消息以提高效率。
			n := len(c.send)
			writeError := false
			for i := 0; i < n; i++ {
				_, err = w.Write(newline) // 消息间添加换行符
				if err != nil {
					global.Logger.Error(fmt.Sprintf("Failed to write newline for %s: %v", c.uuid, err))
					writeError = true
					break
				}
				_, err = w.Write(<-c.send) // 写入下一条排队的消息
				if err != nil {
					global.Logger.Error(fmt.Sprintf("Failed to write queued message for %s: %v", c.uuid, err))
					writeError = true
					break
				}
			}

			// 关闭写入器，将所有数据刷新到底层连接。
			if err := w.Close(); err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to close writer for %s: %v", c.uuid, err))
				return // 关闭写入器失败，通常意味着连接已关闭
			}

			// 如果在写入排队消息时发生错误，则退出 writePump。
			if writeError {
				return
			}

		case <-ticker.C:
			// 定时器触发，发送 Ping 消息。
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to set write deadline for ping for %s: %v", c.uuid, err))
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to send ping for %s: %v", c.uuid, err))
				return // 发送 Ping 失败，通常意味着连接已关闭
			}
		}
	}
}

// ServeWs 处理 WebSocket 连接请求的 HTTP 处理器。
// 负责升级连接、创建 Client、注册到 Hub 并启动读写 goroutine。
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	uid := r.URL.Query().Get("uid")

	// 如果请求中没有 UID，则生成一个新的。
	if uid == "" {
		var err error
		uid, err = encryption.GenerateUID("u_")
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to generate UID: %v", err))
			http.Error(w, "Failed to generate user ID", http.StatusInternalServerError)
			return
		}
	}
	// 验证 UID 格式。
	if ok, err := encryption.ValidateUID(uid, "u_"); err != nil || !ok {
		global.Logger.Warn(fmt.Sprintf("Invalid UID provided or generated: %s, validation error: %v", uid, err))
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 升级 HTTP 连接到 WebSocket。
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for potential user %s: %v", uid, err) // Upgrade 会处理 HTTP 响应
		return
	}

	// 创建 Client 实例，缓冲区大小为 256。
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), uuid: uid, username: username}
	// 注册客户端到 Hub。
	client.hub.register <- client

	// 创建并发送欢迎消息。
	welcomePayload := map[string]string{"username": username, "uid": uid}
	welcomeMessage, err := message_type.NewSystemMessage(welcomePayload).SerializeWithArgs()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Failed to serialize welcome message for %s: %v", client.uuid, err))
	} else {
		client.send <- welcomeMessage // 将欢迎消息放入发送通道
	}

	// 启动后台 goroutine 处理读写。
	go client.writePump()
	go client.readPump()
}
