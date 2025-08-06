// internal/interfaces/websocket/hub_interface.go
package websocket

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type HubInterface interface {
	// Client management
	RegisterClient(conn *websocket.Conn)
	UnregisterClient(conn *websocket.Conn)

	// Broadcasting
	Broadcast(msg []byte)

	// Lifecycle
	Run()
	Stop() error
	IsRunning() bool

	// Client info
	GetClientCount() int
	GetClients() map[*websocket.Conn]bool
}

type HubStats interface {
	GetActiveConnections() int
	GetTotalMessagesSent() int64
	GetTotalMessagesReceived() int64
	GetUptime() time.Duration
}
