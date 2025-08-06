package websockets

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	running    bool
	stopChan   chan struct{}
	stats      *HubStats
}

type HubStats struct {
	activeConnections     int
	totalMessagesSent     int64
	totalMessagesReceived int64
	startTime             time.Time
	mu                    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		stopChan:   make(chan struct{}),
		stats: &HubStats{
			startTime: time.Now(),
		},
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
	h.running = true
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.stats.activeConnections++
			h.mu.Unlock()
			fmt.Printf("Client joined the room: %v\n", conn.RemoteAddr())

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				fmt.Printf("Client left the room: %v\n", conn.RemoteAddr())
				delete(h.clients, conn)
				h.stats.activeConnections--
				conn.Close()
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			h.stats.totalMessagesSent++
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("Error writing message to client:", err)
					conn.Close()
					delete(h.clients, conn)
					h.stats.activeConnections--
				}
			}
			h.mu.Unlock()

		case <-h.stopChan:
			h.running = false
			return
		}
	}
}

func (h *Hub) Stop() error {
	if !h.running {
		return nil
	}

	close(h.stopChan)
	return nil
}

func (h *Hub) IsRunning() bool {
	return h.running
}

func (h *Hub) GetClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

func (h *Hub) GetClients() map[*websocket.Conn]bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Return a copy to avoid race conditions
	clients := make(map[*websocket.Conn]bool)
	for conn, active := range h.clients {
		clients[conn] = active
	}
	return clients
}

func (h *Hub) GetActiveConnections() int {
	h.stats.mu.RLock()
	defer h.stats.mu.RUnlock()
	return h.stats.activeConnections
}

func (h *Hub) GetTotalMessagesSent() int64 {
	h.stats.mu.RLock()
	defer h.stats.mu.RUnlock()
	return h.stats.totalMessagesSent
}

func (h *Hub) GetTotalMessagesReceived() int64 {
	h.stats.mu.RLock()
	defer h.stats.mu.RUnlock()
	return h.stats.totalMessagesReceived
}

func (h *Hub) GetUptime() time.Duration {
	h.stats.mu.RLock()
	defer h.stats.mu.RUnlock()
	return time.Since(h.stats.startTime)
}
