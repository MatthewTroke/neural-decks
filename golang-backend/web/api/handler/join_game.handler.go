package handler

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
)

type JoinGamePayload struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type JoinGameHandler struct {
	Payload JoinGamePayload
	GameSvc domain.GameService
	Claim   *domain.CustomClaim
	Hub     *domain.Hub
}

func NewJoinGameHandler(payload JoinGamePayload, gameSvc domain.GameService, claim *domain.CustomClaim, hub *domain.Hub) *JoinGameHandler {
	return &JoinGameHandler{
		Payload: payload,
		GameSvc: gameSvc,
		Claim:   claim,
		Hub:     hub,
	}
}

func (h *JoinGameHandler) Validate() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

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
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("%s validation failed, could not find game by payload's game id: %w", domain.ContinueRound, err)
	}

	game.Lock()
	defer game.Unlock()

	player := domain.NewPlayer(h.Claim)

	if game.Players == nil {
		return fmt.Errorf("could not join game ID %s, game players is nil", h.Payload.GameID)
	}

	game.Players = append(game.Players, player)

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.JoinGame, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
