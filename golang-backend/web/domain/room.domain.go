package domain

import (
	"fmt"
	"sync"
)

type RoomManager struct {
	rooms map[string]*Hub
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Hub),
	}
}

func (rm *RoomManager) GetRoom(roomID string) *Hub {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; !exists {
		fmt.Println("Creating new game room:", roomID)
		hub := NewHub()
		rm.rooms[roomID] = hub
		go hub.Run()
	}

	return rm.rooms[roomID]
}

func (rm *RoomManager) RemoveRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if hub, exists := rm.rooms[roomID]; exists {
		fmt.Println("Removing empty room:", roomID)
		close(hub.broadcast)
		delete(rm.rooms, roomID)
	}
}
