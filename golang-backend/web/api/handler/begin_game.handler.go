package handler

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
)

type BeginGamePayload struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type BeginGameHandler struct {
	Payload BeginGamePayload
	GameSvc domain.GameService
	Claim   *domain.CustomClaim
	Hub     *domain.Hub
}

func NewBeginGameHandler(payload BeginGamePayload, gameSvc domain.GameService, claim *domain.CustomClaim, hub *domain.Hub) *BeginGameHandler {
	return &BeginGameHandler{
		Payload: payload,
		GameSvc: gameSvc,
		Claim:   claim,
		Hub:     hub,
	}
}

func (h *BeginGameHandler) Validate() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to validate inbound %s event: %w", domain.BeginGame, err)
	}

	if game.Players == nil {
		return fmt.Errorf("unable to begin game players list is nil")
	}

	if len(game.Players) < 2 {
		return fmt.Errorf("unable to begin game players list, not enough players")
	}

	return nil
}

func (h *BeginGameHandler) Handle() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	game.Lock()
	defer game.Unlock()

	game.Collection.Shuffle()
	game.Players[0].SetIsCardCzar(true)

	for i := range game.Players {
		game.Players[i].Deck = game.Collection.DrawCards(7, "White")
	}

	game.BlackCard = game.Collection.DrawCards(1, domain.Black)[0]

	game.SetRoundStatus(domain.PlayersPickingCard)
	game.SetStatus(domain.InProgress)

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
