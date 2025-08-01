package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"errors"
	"fmt"
	"time"
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

	// Check if anyone has won the game
	var winner *domain.Player
	for _, player := range newGame.Players {
		if player.Score >= newGame.WinnerCount {
			winner = player
			break
		}
	}

	if winner != nil {
		// Game is over, someone won! Create and apply game winner event
		gameWinnerEvent, err := h.EventService.CreateGameEvent(
			h.Payload.GameID,
			domain.EventGameWinner,
			domain.NewGameEventPayloadGameWinner(h.Payload.GameID, winner.UserID, winner.Score),
		)

		if err != nil {
			return fmt.Errorf("failed to create game winner event: %w", err)
		}

		// Apply the game winner event to the new game state
		if err := newGame.ApplyEvent(gameWinnerEvent); err != nil {
			return fmt.Errorf("failed to apply game winner event: %w", err)
		}

		// Persist the game winner event
		if err := h.EventService.AppendEvent(gameWinnerEvent); err != nil {
			return fmt.Errorf("failed to persist game winner event: %w", err)
		}

		message := domain.NewWebSocketMessage(domain.GameUpdate, newGame)
		chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, fmt.Sprintf("🎉 %s has won the game with %d points! 🎉", winner.Name, winner.Score))

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

		// Stop the timer immediately to prevent logging spam
		h.GameStateService.StopGameTimer(h.Payload.GameID)

		// Schedule cleanup of the finished game after 30 seconds
		go func() {
			time.Sleep(30 * time.Second)
			h.GameStateService.CleanupGame(h.Payload.GameID)
		}()
	} else {
		// Normal round continuation
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

		// Only reset timer if game is still in progress
		h.GameStateService.ResetAutoContinueTimer(h.Payload.GameID)
	}

	return nil
}
