package websocket

import (
	"fmt"

	"github.com/qianmianyao/parchment-server/internal/services/chat"
	"github.com/qianmianyao/parchment-server/pkg/global"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	chatCreate *chat.Create
	chatUpdate *chat.Update
	chatFind   *chat.Find
}

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

// Run starts the hub's main loop
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
		if err := h.chatCreate.User(client.username, client.uuid); err != nil {
			return
		}
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
		global.Logger.Debug(fmt.Sprintf("client %v disconnected", client))
		// 从在线用户列表中删除
		delete(usersClients, client.uuid)
		delete(h.clients, client)
		close(client.send)
	}
}

// Broadcast sends a message_type to all connected clients
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

// SendToSpecificClient 发送消息给指定客户端
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
