package websocket

import (
	"fmt"
	"github.com/qianmianyao/parchment-server/pkg/global"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// run starts the hub's main loop
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			global.Logger.Debug(fmt.Sprintf("client %v connected", client))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				global.Logger.Debug(fmt.Sprintf("client %v disconnected", client))
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
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
	}
}
