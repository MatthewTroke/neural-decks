package ws

import (
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn   *websocket.Conn
	RoomID string
	Send   chan []byte
}

func (c *Client) writePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	_ = c.Conn.Close()
}
