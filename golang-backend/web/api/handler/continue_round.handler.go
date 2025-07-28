package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
)

type ContinueRoundPayload struct {
	GameID string `json:"game_id"`
}

type ContinueRoundHandler struct {
	Payload          request.GameEventPayloadGameRoundContinuedRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewContinueRoundHandler(payload request.GameEventPayloadGameRoundContinuedRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *ContinueRoundHandler {
	return &ContinueRoundHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *ContinueRoundHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", domain.ContinueRound, err)
	}

	_, err = game.FindPlayerByUserId(h.Claim.UserID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find player by user id: %w", domain.ContinueRound, err)
	}

	return nil
}

func (h *ContinueRoundHandler) Handle() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventRoundContinued,
		domain.NewGameEventPayloadGameRoundContinued(h.Payload.GameID, h.Payload.UserID),
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	newGame := game.Clone()

	if err := newGame.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}

	if err := h.EventService.AppendEvent(event); err != nil {
		return fmt.Errorf("failed to persist event: %w", err)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, newGame)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.ContinueRound, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
