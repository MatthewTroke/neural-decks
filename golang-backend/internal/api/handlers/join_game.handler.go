//go:build exclude

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

type JoinGamePayload struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type JoinGameHandler struct {
	Payload          request.GameEventPayloadJoinedGameRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *entities.CustomClaim
	Hub              *websockets.Hub
}

func NewJoinGameHandler(payload request.GameEventPayloadJoinedGameRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *entities.CustomClaim, hub *websockets.Hub) *JoinGameHandler {
	return &JoinGameHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *JoinGameHandler) Validate() error {
	game, err := h.EventService.BuildGameByGameId(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", aggregates.ContinueRound, err)
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
	currentGame, err := h.EventService.BuildGameByGameId(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", aggregates.ContinueRound, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		events.EventJoinedGame,
		aggregates.NewGameEventPayloadJoinedGame(h.Payload.GameID, h.Payload.UserID, h.Claim),
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

	log := fmt.Sprintf("%s has joined the game.", h.Claim.Name)

	message := websockets.NewWebSocketMessage(aggregates.GameUpdate, newGame)
	chatMessage := websockets.NewWebSocketMessage(aggregates.ChatMessage, log)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.JoinGame, err)
	}

	jsonChatMessage, err := json.Marshal(chatMessage)

	if err != nil {
		return fmt.Errorf("unable to marshal chat message: %w", err)
	}

	h.Hub.Broadcast(jsonMessage)
	h.Hub.Broadcast(jsonChatMessage)

	return nil
}
