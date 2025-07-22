package handler

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
)

type PlayCardPayload struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type PlayCardHandler struct {
	Payload PlayCardPayload
	GameSvc domain.GameService
	Claim   *domain.CustomClaim
	Hub     *domain.Hub
}

func NewPlayCardHandler(payload PlayCardPayload, gameSvc domain.GameService, claim *domain.CustomClaim, hub *domain.Hub) *PlayCardHandler {
	return &PlayCardHandler{
		Payload: payload,
		GameSvc: gameSvc,
		Claim:   claim,
		Hub:     hub,
	}
}

func (h *PlayCardHandler) Validate() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	player, err := game.FindPlayerByUserId(h.Claim.UserID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	_, err = game.FindCardByPlayerId(player.UserID, h.Payload.CardID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	if player.HasAlreadyPlayedWhiteCard() {
		return fmt.Errorf("unable to handle inbound %s event, player has already played a white card", domain.PlayCard)
	}

	if player.IsCardCzar {
		return fmt.Errorf("unable to handle inbound %s event, player is currently a card czar", domain.PlayCard)
	}

	return nil
}

func (h *PlayCardHandler) Handle() error {
	game, err := h.GameSvc.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	player, _ := game.FindPlayerByUserId(h.Claim.UserID)
	card, _ := game.FindCardByPlayerId(player.UserID, h.Payload.CardID)

	err = player.RemoveCardFromDeck(card.ID)

	if err != nil {
		return fmt.Errorf("unable to play white card: %w", err)
	}

	err = player.SetCardAsPlacedCard(card)

	if err != nil {
		return fmt.Errorf("unable to play white card: %w", err)
	}

	err = game.AddWhiteCardToGameBoard(card)

	if err != nil {
		return fmt.Errorf("unable to play white card: %w", err)
	}

	hasPlayersPlayedWhiteCard, err := game.HasAllPlayersPlayedWhiteCard()

	if err != nil {
		return fmt.Errorf("unable to play white card: %w", err)
	}

	if hasPlayersPlayedWhiteCard {
		game.SetRoundStatus(domain.CardCzarPickingWinningCard)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	h.Hub.Broadcast(jsonMessage)

	return nil
}
