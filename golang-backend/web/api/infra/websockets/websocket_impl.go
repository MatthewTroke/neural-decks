package websockets

type WebsocketService interface {
	TriggerEvent(channel string, event string, data interface{}) error
}

type WebSocketHandler interface {
	Handle() error
}

type WebSocketGameEvent[M any, P any] struct {
	Type    M `json:"type"`
	Payload P `json:"payload"`
}

func NewWebSocketMessage[M, P any](messageType M, payload P) WebSocketGameEvent[M, P] {
	return WebSocketGameEvent[M, P]{
		Type:    messageType,
		Payload: payload,
	}
}
