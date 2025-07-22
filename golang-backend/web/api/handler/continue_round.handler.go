package handler

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
)

type ContinueRoundPayload struct {
	GameID string `json:"game_id"`
}

type ContinueRoundHandler struct {
	Payload ContinueRoundPayload
	GameSvc domain.GameService
	Claim   *domain.CustomClaim
	Hub     *domain.Hub
}

func NewContinueRoundHandler(payload ContinueRoundPayload, gameSvc domain.GameService, claim *domain.CustomClaim, hub *domain.Hub) *ContinueRoundHandler {
	return &ContinueRoundHandler{
		Payload: payload,
		GameSvc: gameSvc,
		Claim:   claim,
		Hub:     hub,
	}
}

func (h *ContinueRoundHandler) Validate() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

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
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	game.Lock()
	defer game.Unlock()

	// game.RemovePlacedCards()
	// game.DrawWhiteCardsFromAllPlayersWhoPlayed()

	for _, player := range game.Players {
		if player.IsCardCzar {
			continue
		}

		player.RemovePlacedCard()
		player.Deck = append(player.Deck, game.Collection.DrawCards(1, domain.White)...)
	}

	currentCardCzar, err := game.FindCurrentCardCzar()

	if err != nil {
		return fmt.Errorf("could not continue round: %w", err)
	}

	currentCardCzar.SetIsCardCzar(false)
	currentCardCzar.SetWasCardCzar(true)

	game.SetRoundWinner(nil)
	game.ClearBoard()

	err = game.PickNewCardCzar()

	if err != nil {
		return fmt.Errorf("could not pick new card czar: %w", err)
	}

	err = game.PickNewBlackCard()

	if err != nil {
		return fmt.Errorf("could not pick new black card: %w", err)
	}

	game.IncrementGameRound()
	game.SetRoundStatus(domain.PlayersPickingCard)

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.ContinueRound, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
