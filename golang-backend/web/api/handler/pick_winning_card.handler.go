package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"errors"
	"fmt"
)

type PickWinningCardPayload struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type PickWinningCardHandler struct {
	Payload          request.GameEventPayloadCardCzarChoseWinningCardRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
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

func NewPickWinningCardHandler(payload request.GameEventPayloadCardCzarChoseWinningCardRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *PickWinningCardHandler {
	return &PickWinningCardHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *PickWinningCardHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

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
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PickWinningCard, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventCardCzarChoseWinningCard,
		domain.NewGameEventPayloadCardCzarChoseWinningCard(h.Payload.GameID, h.Payload.CardID),
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
	chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, "Card czar has chosen a winning card.")

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return errors.New("failed to marshal message")
	}

	jsonChatMessage, err := json.Marshal(chatMessage)

	if err != nil {
		return fmt.Errorf("unable to marshal chat message: %w", err)
	}

	h.Hub.Broadcast(jsonMessage)
	h.Hub.Broadcast(jsonChatMessage)

	return nil
}
