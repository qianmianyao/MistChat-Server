package websocket

import (
	"fmt"

	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/pkg/global"
)

// Hub 负责管理 WebSocket 客户端连接、注册、注销以及消息广播。
type Hub struct {
	// clients 存储所有已连接的客户端。
	clients map[*Client]bool
	// broadcast 通道用于接收需要广播给所有客户端的消息。
	broadcast chan []byte
	// register 通道用于接收新客户端的注册请求。
	register chan *Client
	// unregister 通道用于接收客户端的注销请求。
	unregister chan *Client
	// chatCreate 用于处理聊天相关的创建操作。
	chatCreate *chat.Create
	// chatUpdate 用于处理聊天相关的更新操作。
	chatUpdate *chat.Update
	// chatFind 用于处理聊天相关的查找操作。
	chatFind *chat.Find
}

// NewHub 创建并返回一个新的 Hub 实例。
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		chatCreate: chat.NewCreate(),
		chatUpdate: chat.NewUpdate(),
		chatFind:   chat.NewFind(),
	}
}

var usersClients = make(map[string]*Client)

// Run 启动 Hub 的主事件循环，监听并处理客户端注册、注销和消息广播事件。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clientRegister(client)
		case client := <-h.unregister:
			h.clientUnregister(client)
		case message := <-h.broadcast:
			h.Broadcast(message)
		}
	}
}

// clientRegister registers a new client
func (h *Hub) clientRegister(client *Client) {
	r := h.chatFind.IsUserExist(client.uuid)
	switch r {
	case chat.UserExist:
		if err := h.chatUpdate.UserOnlineStatus(client.uuid, true); err != nil {
			return
		}
	case chat.UserNotExist:
		return
	}
	h.clients[client] = true
	usersClients[client.uuid] = client
	global.Logger.Debug(fmt.Sprintf("client %v connected", client))
}

// clientUnregister unregisters a client
func (h *Hub) clientUnregister(client *Client) {
	if _, ok := h.clients[client]; ok {

		if r := h.chatFind.IsUserExist(client.uuid); r == chat.UserExist {
			if err := h.chatUpdate.UserOnlineStatus(client.uuid, false); err != nil {
				return
			}
		}
		// 从在线用户列表中删除
		delete(usersClients, client.uuid)
		delete(h.clients, client)
	}
}

// Broadcast 将消息发送给所有连接的客户端。
func (h *Hub) Broadcast(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
			global.Logger.Debug(fmt.Sprintf("send to client: %v", client))
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// SendToSpecificClient 将消息发送给指定房间内除发送者外的所有其他客户端。
// uuid: 发送者客户端的UUID。
// roomUUID: 目标房间的UUID。
// message: 要发送的消息内容。
func (h *Hub) SendToSpecificClient(uuid, roomUUID string, message []byte) {
	// 获取房间内所有的用户
	users := h.chatFind.AllUsersInTheRoom(roomUUID)

	roomstatus := h.chatFind.IsTheUserIsInTheRoom(uuid, roomUUID)
	if roomstatus == chat.NotInRoom {
		global.Logger.Debug(fmt.Sprintf("用户 %s 不在房间 %s 内", uuid, roomUUID))
		return
	}

	// 获取房间内所有的在线用户
	var clients []*Client
	for _, uid := range users {
		if client, ok := usersClients[uid]; ok {
			// 不对自己发送消息
			if client.uuid == uuid {
				continue
			}
			clients = append(clients, client)
		}
	}
	for _, client := range clients {
		select {
		case client.send <- message:
			global.Logger.Debug(fmt.Sprintf("send to specific client: %v", client))
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
