package handlers

import (
	"cardgame/domain/aggregates"
	"cardgame/domain/entities"
	"cardgame/domain/events"
	"cardgame/domain/services"
	"cardgame/http/request"
	"cardgame/infra/websockets"
	"encoding/json"
	"fmt"
)

type EmojiClickedHandler struct {
	Payload          request.GameEventPayloadEmojiClickedRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *entities.CustomClaim
	Hub              *websockets.Hub
}

func NewEmojiClickedHandler(payload request.GameEventPayloadEmojiClickedRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *entities.CustomClaim, hub *websockets.Hub) *EmojiClickedHandler {
	return &EmojiClickedHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *EmojiClickedHandler) Validate() error {
	game, err := h.EventService.BuildGameByGameId(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to validate inbound %s event: %w", events.EventEmojiClicked, err)
	}

	// Check if the user is a player in the game
	_, err = game.FindPlayerByUserId(h.Payload.UserID)
	if err != nil {
		return fmt.Errorf("user %s is not a player in game %s", h.Payload.UserID, h.Payload.GameID)
	}

	// Validate emoji is not empty
	if h.Payload.Emoji == "" {
		return fmt.Errorf("emoji cannot be empty")
	}

	return nil
}

func (h *EmojiClickedHandler) Handle() error {
	message := websockets.NewWebSocketMessage(aggregates.EmojiClickedMessage, map[string]interface{}{
		"user_id": h.Payload.UserID,
		"emoji":   h.Payload.Emoji,
		"game_id": h.Payload.GameID,
	})

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("unable to marshal emoji clicked message: %w", err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
