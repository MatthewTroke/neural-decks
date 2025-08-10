package interfaces

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type Hub interface {
	Broadcast(msg []byte)

	Run()
	Stop() error
	IsRunning() bool

	GetClientCount() int
	GetClients() map[*websocket.Conn]bool

	GetActiveConnections() int
	GetTotalMessagesSent() int64
	GetTotalMessagesReceived() int64
	GetUptime() time.Duration
}
