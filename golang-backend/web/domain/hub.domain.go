package domain

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	Clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) RegisterClient(conn *websocket.Conn) {
	h.register <- conn
}

func (h *Hub) UnregisterClient(conn *websocket.Conn) {
	h.unregister <- conn
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.Clients[conn] = true
			h.mu.Unlock()
			fmt.Printf("Client joined the room: %v\n", conn.RemoteAddr())

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.Clients[conn]; ok {
				fmt.Printf("Client left the room: %v\n", conn.RemoteAddr())
				delete(h.Clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for conn := range h.Clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("Error writing message to client:", err)
					conn.Close()
					delete(h.Clients, conn)
				}
			}
			h.mu.Unlock()
		}
	}
}
