package ws

import "encoding/json"

type Publisher struct {
	Hub *Hub
}

func NewPublisher(h *Hub) *Publisher { return &Publisher{Hub: h} }

func (p *Publisher) PublishToRoom(roomID string, eventType string, payload any) error {
	envelope := map[string]any{
		"type":    eventType,
		"payload": payload,
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	p.Hub.Broadcast(roomID, b)
	return nil
}
