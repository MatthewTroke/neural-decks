package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ContinueRoundPayload struct {
	GameID string `json:"game_id"`
}

type ContinueRoundHandler struct {
	Payload          request.GameEventPayloadGameRoundContinuedRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewContinueRoundHandler(payload request.GameEventPayloadGameRoundContinuedRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *ContinueRoundHandler {
	return &ContinueRoundHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *ContinueRoundHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

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
	game, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	// Get used cards from Redis for checking
	usedCardIDs, err := h.EventService.GetUsedCards(h.Payload.GameID)
	if err != nil {
		return fmt.Errorf("failed to get used cards from Redis: %w", err)
	}

	// Convert to map for efficient lookup
	usedCardsMap := make(map[string]bool)
	for _, cardID := range usedCardIDs {
		usedCardsMap[cardID] = true
	}

	// Count how many non-card-czar players need cards
	playersNeedingCards := 0
	for _, player := range game.Players {
		if !player.IsCardCzar {
			playersNeedingCards++
		}
	}

	// Count available unused white cards
	availableWhiteCards := 0
	for _, card := range game.Collection.Cards {
		if card.Type == "White" && !usedCardsMap[card.ID] {
			availableWhiteCards++
		}
	}

	// If we don't have enough cards, create a shuffle event first
	if availableWhiteCards < playersNeedingCards {
		// Clear used cards in Redis
		if err := h.EventService.ClearUsedCards(h.Payload.GameID); err != nil {
			return fmt.Errorf("failed to clear used cards: %w", err)
		}

		// Create shuffle event
		shuffleEvent, err := h.EventService.CreateGameEvent(
			h.Payload.GameID,
			domain.EventShuffle,
			domain.NewGameEventPayloadShuffle(h.Payload.GameID, time.Now().UnixNano(), uuid.New().String()),
		)

		if err != nil {
			return fmt.Errorf("failed to create shuffle event: %w", err)
		}

		chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, "No more available white cards. Re-shuffling deck...")

		jsonChatMessage, err := json.Marshal(chatMessage)

		if err != nil {
			return fmt.Errorf("unable to marshal chat message: %w", err)
		}

		h.Hub.Broadcast(jsonChatMessage)

		if err != nil {
			return fmt.Errorf("failed to create shuffle event: %w", err)
		}

		// Apply shuffle event
		if err := game.ApplyEvent(shuffleEvent); err != nil {
			return fmt.Errorf("failed to apply shuffle event: %w", err)
		}

		// Persist shuffle event
		if err := h.EventService.AppendEvent(shuffleEvent); err != nil {
			return fmt.Errorf("failed to persist shuffle event: %w", err)
		}

		// Reset used cards map since we just shuffled
		usedCardsMap = make(map[string]bool)
	}

	// Determine which cards to give to each player
	playerCards := make(map[string]string)
	newlyUsedCards := []string{}

	for _, player := range game.Players {
		if player.IsCardCzar {
			continue
		}

		// Find an unused white card for this player
		for _, card := range game.Collection.Cards {
			if card.Type == "White" && !usedCardsMap[card.ID] {
				playerCards[player.UserID] = card.ID
				newlyUsedCards = append(newlyUsedCards, card.ID)
				break
			}
		}
	}

	newGame := game.Clone()

	// Check if we have enough unused black cards
	availableBlackCards := 0
	for _, card := range game.Collection.Cards {
		if card.Type == "Black" && !usedCardsMap[card.ID] {
			availableBlackCards++
		}
	}

	// If we don't have enough black cards, create another shuffle event
	if availableBlackCards < 1 {
		// Clear used cards in Redis
		if err := h.EventService.ClearUsedCards(h.Payload.GameID); err != nil {
			return fmt.Errorf("failed to clear used cards: %w", err)
		}

		// Create shuffle event
		shuffleEvent, err := h.EventService.CreateGameEvent(
			h.Payload.GameID,
			domain.EventShuffle,
			domain.NewGameEventPayloadShuffle(h.Payload.GameID, time.Now().UnixNano(), uuid.New().String()),
		)

		if err != nil {
			return fmt.Errorf("failed to create shuffle event: %w", err)
		}

		chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, "No more available black cards. Re-shuffling deck...")

		jsonChatMessage, err := json.Marshal(chatMessage)

		if err != nil {
			return fmt.Errorf("unable to marshal chat message: %w", err)
		}

		h.Hub.Broadcast(jsonChatMessage)

		if err != nil {
			return fmt.Errorf("failed to create shuffle event: %w", err)
		}

		// Apply shuffle event
		if err := game.ApplyEvent(shuffleEvent); err != nil {
			return fmt.Errorf("failed to apply shuffle event: %w", err)
		}

		// Persist shuffle event
		if err := h.EventService.AppendEvent(shuffleEvent); err != nil {
			return fmt.Errorf("failed to persist shuffle event: %w", err)
		}

		// Reset used cards map since we just shuffled
		usedCardsMap = make(map[string]bool)
	}

	// Find an unused black card
	var blackCardID string
	for _, card := range game.Collection.Cards {
		if card.Type == "Black" && !usedCardsMap[card.ID] {
			blackCardID = card.ID
			newlyUsedCards = append(newlyUsedCards, card.ID)
			break
		}
	}

	// Create event with the specific cards for each player
	event, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventRoundContinued,
		domain.NewGameEventPayloadGameRoundContinuedWithCards(h.Payload.GameID, h.Payload.UserID, playerCards, blackCardID),
	)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	if err := newGame.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}

	// Batch add newly used cards to Redis
	if len(newlyUsedCards) > 0 {
		if err := h.EventService.AddUsedCards(h.Payload.GameID, newlyUsedCards); err != nil {
			return fmt.Errorf("failed to add used cards to Redis: %w", err)
		}
	}

	if err := h.EventService.AppendEvent(event); err != nil {
		return fmt.Errorf("failed to persist event: %w", err)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, newGame)
	chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, "Round has continued.")

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to marshal game update: %w", err)
	}

	jsonChatMessage, err := json.Marshal(chatMessage)

	if err != nil {
		return fmt.Errorf("unable to marshal chat message: %w", err)
	}

	h.Hub.Broadcast(jsonMessage)
	h.Hub.Broadcast(jsonChatMessage)

	h.GameStateService.ResetAutoContinueTimer(h.Payload.GameID)

	return nil
}
