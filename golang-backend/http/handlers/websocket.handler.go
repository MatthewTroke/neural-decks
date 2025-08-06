package handlers

type WebSocketHandler interface {
	Validate() error
	Handle() error
}
