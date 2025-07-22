package handler

type WebSocketHandler interface {
	Validate() error
	Handle() error
}
