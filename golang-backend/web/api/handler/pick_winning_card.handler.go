package handler

import (
	"cardgame/domain"
	"encoding/json"
	"errors"
	"fmt"
)

type PickWinningCardPayload struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type PickWinningCardHandler struct {
	Payload PickWinningCardPayload
	GameSvc domain.GameService
	Claim   *domain.CustomClaim
	Hub     *domain.Hub
}

type InvalidUserActionError struct {
	Message string
	Type    string
}

func (e *InvalidUserActionError) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *InvalidUserActionError) Error() string {
	return e.Message
}

func NewPickWinningCardHandler(payload PickWinningCardPayload, gameSvc domain.GameService, claim *domain.CustomClaim, hub *domain.Hub) *PickWinningCardHandler {
	return &PickWinningCardHandler{
		Payload: payload,
		GameSvc: gameSvc,
		Claim:   claim,
		Hub:     hub,
	}
}

func (h *PickWinningCardHandler) Validate() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PickWinningCard, err)
	}

	if game.Status != domain.InProgress {
		return fmt.Errorf("unable to handle inbound %s event, game status is not InProgress", domain.PickWinningCard)
	}

	if game.RoundStatus != domain.CardCzarPickingWinningCard {
		return fmt.Errorf("unable to handle inbound %s event, round status is not in CardCzarPickingWinningCard", domain.PickWinningCard)
	}

	hasPlayersPlayedWhiteCard, err := game.HasAllPlayersPlayedWhiteCard()

	if err != nil {
		return fmt.Errorf("unable to handle %s event: %w", domain.PickWinningCard, err)
	}

	if !hasPlayersPlayedWhiteCard {
		return fmt.Errorf("unable to handle inbound %s event, not all players have played a white card", domain.PickWinningCard)
	}

	player, err := game.FindPlayerByUserId(h.Claim.UserID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event, unable to find player with given user id: %w", domain.PickWinningCard, err)
	}

	// if !player.IsCardCzar {
	// 	return fmt.Errorf("unable to handle inbound %s event, player picking winning card is not a card czar: %w", domain.PickWinningCard, err)
	// }

	if !player.IsCardCzar {
		return &InvalidUserActionError{
			Message: "Sorry! you cant do that",
			Type:    "PLAYER_IS_NOT_CARD_CZAR",
		}
	}

	winningCard, err := game.FindWhiteCardByCardId(h.Payload.CardID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event, unable to find winning card: %w", domain.PickWinningCard, err)
	}

	_, err = game.FindWhiteCardOwner(winningCard)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event, unable to find winning card owner: %w", domain.PickWinningCard, err)
	}

	return nil
}

func (h *PickWinningCardHandler) Handle() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	winningCard, _ := game.FindWhiteCardByCardId(h.Payload.CardID)
	winner, _ := game.FindWhiteCardOwner(winningCard)

	winner.IncrementScore()
	game.SetRoundStatus(domain.CardCzarChoseWinningCard)
	game.SetRoundWinner(winner)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event, unable to choose white card winner: %w", domain.PickWinningCard, err)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return errors.New("failed to marshal message")
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
