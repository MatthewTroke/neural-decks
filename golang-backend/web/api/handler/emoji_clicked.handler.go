package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
)

type EmojiClickedHandler struct {
	Payload          request.GameEventPayloadEmojiClickedRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewEmojiClickedHandler(payload request.GameEventPayloadEmojiClickedRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *EmojiClickedHandler {
	return &EmojiClickedHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *EmojiClickedHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to validate inbound %s event: %w", domain.EventEmojiClicked, err)
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
	message := domain.NewWebSocketMessage(domain.EmojiClickedMessage, map[string]interface{}{
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
