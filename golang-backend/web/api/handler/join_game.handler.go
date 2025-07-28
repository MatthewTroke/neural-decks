package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
)

type JoinGamePayload struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type JoinGameHandler struct {
	Payload          request.GameEventPayloadJoinedGameRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewJoinGameHandler(payload request.GameEventPayloadJoinedGameRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *JoinGameHandler {
	return &JoinGameHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *JoinGameHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", domain.ContinueRound, err)
	}

	if len(game.Players) >= game.MaxPlayerCount {
		return fmt.Errorf("could not join game ID %s, max player count reached", game.ID)
	}

	for _, player := range game.Players {
		if player.UserID == h.Claim.UserID {
			return fmt.Errorf("could not join game ID %s, player with id already %s in game", game.ID, player.UserID)
		}
	}

	return nil
}

func (h *JoinGameHandler) Handle() error {
	currentGame, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", domain.ContinueRound, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventJoinedGame,
		domain.NewGameEventPayloadJoinedGame(h.Payload.GameID, h.Payload.UserID, h.Claim),
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	newGame := currentGame.Clone()

	if err := newGame.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}

	if err := h.EventService.AppendEvent(event); err != nil {
		return fmt.Errorf("failed to persist event: %w", err)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, newGame)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.JoinGame, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
