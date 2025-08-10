package ws

import "sync"

type Hub struct {
	mu       sync.RWMutex
	rooms    map[string]map[*Client]struct{}
	register chan *Client
	unreg    chan *Client
	bcast    chan broadcast
}

type broadcast struct {
	roomID string
	data   []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:    map[string]map[*Client]struct{}{},
		register: make(chan *Client, 64),
		unreg:    make(chan *Client, 64),
		bcast:    make(chan broadcast, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[c.RoomID]; !ok {
				h.rooms[c.RoomID] = map[*Client]struct{}{}
			}
			h.rooms[c.RoomID][c] = struct{}{}
			h.mu.Unlock()
			go c.writePump()

		case c := <-h.unreg:
			h.mu.Lock()
			if clients, ok := h.rooms[c.RoomID]; ok {
				delete(clients, c)
				if len(clients) == 0 {
					delete(h.rooms, c.RoomID)
				}
			}
			h.mu.Unlock()
			close(c.Send)

		case m := <-h.bcast:
			h.mu.RLock()
			if clients, ok := h.rooms[m.roomID]; ok {
				for c := range clients {
					select {
					case c.Send <- m.data:
					default:
						// slow client, drop connection
						go func(c *Client) { h.unreg <- c }(c)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Join(roomID string, c *Client) { h.register <- c }
func (h *Hub) Leave(c *Client)               { h.unreg <- c }
func (h *Hub) Broadcast(roomID string, data []byte) {
	h.bcast <- broadcast{roomID: roomID, data: data}
}
