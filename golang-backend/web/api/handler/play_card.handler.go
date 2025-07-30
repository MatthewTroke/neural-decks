package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
)

type PlayCardPayload struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type PlayCardHandler struct {
	Payload          request.GameEventPayloadPlayCardRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewPlayCardHandler(payload request.GameEventPayloadPlayCardRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *PlayCardHandler {
	return &PlayCardHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *PlayCardHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	player, err := game.FindPlayerByUserId(h.Claim.UserID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	_, err = game.FindCardByPlayerId(h.Claim.UserID, h.Payload.CardID)

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
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventCardPlayed,
		domain.NewGameEventPayloadPlayCard(h.Payload.GameID, h.Payload.CardID, h.Claim),
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
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.PlayCard, err)
	}

	h.Hub.Broadcast(jsonMessage)

	hasAllPlayed, err := newGame.HasAllPlayersPlayedWhiteCard()

	if err == nil && hasAllPlayed {
		h.GameStateService.ResetAutoContinueTimer(h.Payload.GameID)
	}

	return nil
}
