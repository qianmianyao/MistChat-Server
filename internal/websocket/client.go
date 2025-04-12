package websocket

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/qianmianyao/parchment-server/internal/websocket/message_type"
	"github.com/qianmianyao/parchment-server/pkg/global"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message_type to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message_type from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message_type size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的请求
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	uuid     string
	username string
}

// readPump reads messages from the websocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Warn(fmt.Sprintf("websocket close: %v", err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}

}

// writePump writes messages to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)

			// Add queued chat messages to the current websocket message_type.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				_, err = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("websocket close: %v", err))
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// user uuid
	uid := r.URL.Query().Get("uid")
	username := r.URL.Query().Get("username")

	if uid == "" {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), uuid: uid, username: username}
	client.hub.register <- client

	welcomeMessage, err := message_type.NewSystemMessage(fmt.Sprintf("Welcome User %v", uid)).Serialize()
	client.send <- welcomeMessage

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
