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

type PlayCardPayload struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type PlayCardHandler struct {
	Payload          request.GameEventPayloadPlayCardRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *entities.CustomClaim
	Hub              *websockets.Hub
}

func NewPlayCardHandler(payload request.GameEventPayloadPlayCardRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *entities.CustomClaim, hub *websockets.Hub) *PlayCardHandler {
	return &PlayCardHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *PlayCardHandler) Validate() error {
	game, err := h.EventService.BuildGameByGameId(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.PlayCard, err)
	}

	player, err := game.FindPlayerByUserId(h.Claim.UserID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.PlayCard, err)
	}

	_, err = game.FindCardByPlayerId(h.Claim.UserID, h.Payload.CardID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.PlayCard, err)
	}

	if player.HasAlreadyPlayedWhiteCard() {
		return fmt.Errorf("unable to handle inbound %s event, player has already played a white card", aggregates.PlayCard)
	}

	if player.IsJudge {
		return fmt.Errorf("unable to handle inbound %s event, player is currently a judge", aggregates.PlayCard)
	}

	return nil
}

func (h *PlayCardHandler) Handle() error {
	game, err := h.EventService.BuildGameByGameId(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.PlayCard, err)
	}

	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		events.EventCardPlayed,
		aggregates.NewGameEventPayloadPlayCard(h.Payload.GameID, h.Payload.CardID, h.Claim),
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

	message := websockets.NewWebSocketMessage(aggregates.GameUpdate, newGame)

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", aggregates.PlayCard, err)
	}

	h.Hub.Broadcast(jsonMessage)

	// Alternative: Use GameCoordinator broadcasting (if GameCoordinator is available)
	// This would send the event directly to all players via their websocket connections
	// instead of broadcasting to the hub
	/*
		if gameCoordinator != nil {
			// Broadcast the card played event to all players
			err := gameCoordinator.BroadcastCardPlayed(h.Payload.GameID, h.Claim.UserID, h.Payload.CardID)
			if err != nil {
				fmt.Printf("Warning: Could not broadcast card played event: %v\n", err)
			}

			// Broadcast game update to all players
			err = gameCoordinator.BroadcastGameUpdate(h.Payload.GameID, newGame)
			if err != nil {
				fmt.Printf("Warning: Could not broadcast game update: %v\n", err)
			}
		}
	*/

	hasAllPlayed, err := newGame.HasAllPlayersPlayedWhiteCard()

	if err == nil && hasAllPlayed {
		h.GameStateService.ResetAutoContinueTimer(h.Payload.GameID)
	}

	return nil
}
