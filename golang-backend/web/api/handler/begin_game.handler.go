package handler

import (
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"encoding/json"
	"fmt"
	"hash/fnv"
)

// hash creates a deterministic hash from a string
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type BeginGamePayload struct {
	GameID string `json:"game_id"`
}

type BeginGameHandler struct {
	Payload          request.GameEventPayloadGameBeginsRequest
	EventService     *services.EventService
	GameStateService *services.GameStateService
	Claim            *domain.CustomClaim
	Hub              *domain.Hub
}

func NewBeginGameHandler(payload request.GameEventPayloadGameBeginsRequest, eventService *services.EventService, gameStateService *services.GameStateService, claim *domain.CustomClaim, hub *domain.Hub) *BeginGameHandler {
	return &BeginGameHandler{
		Payload:          payload,
		EventService:     eventService,
		GameStateService: gameStateService,
		Claim:            claim,
		Hub:              hub,
	}
}

func (h *BeginGameHandler) Validate() error {
	game, err := h.EventService.GetGameById(h.Payload.GameID)

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
	currentGame, err := h.EventService.GetGameById(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("unable to handle inbound %s event: %w", domain.BeginGame, err)
	}

	// 1. Create the game begins event (sets up game state)
	gameBeginsEvent, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventGameBegins,
		domain.NewGameEventPayloadGameBegins(h.Payload.GameID, h.Payload.UserID),
	)
	if err != nil {
		return fmt.Errorf("failed to create game begins event: %w", err)
	}

	// Apply game begins event
	game := currentGame.Clone()

	if err := game.ApplyEvent(gameBeginsEvent); err != nil {
		return fmt.Errorf("failed to apply game begins event: %w", err)
	}

	// Set the first player as card czar
	setCardCzarEvent, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventSetCardCzar,
		domain.NewGameEventPayloadSetCardCzar(h.Payload.GameID, game.Players[0].UserID),
	)

	if err != nil {
		return fmt.Errorf("failed to create set card czar event: %w", err)
	}

	// Apply set card czar event
	if err := game.ApplyEvent(setCardCzarEvent); err != nil {
		return fmt.Errorf("failed to apply set card czar event: %w", err)
	}

	// Persist set card czar event
	if err := h.EventService.AppendEvent(setCardCzarEvent); err != nil {
		return fmt.Errorf("failed to persist set card czar event: %w", err)
	}

	// 2. Create shuffle event with deterministic seed
	// Use a hash of the game ID to ensure deterministic shuffling
	shuffleSeed := int64(hash(h.Payload.GameID))

	shuffleEvent, err := h.EventService.CreateGameEvent(
		h.Payload.GameID,
		domain.EventShuffle,
		domain.NewGameEventPayloadShuffle(h.Payload.GameID, shuffleSeed, fmt.Sprintf("shuffle_%d", shuffleSeed)),
	)

	if err != nil {
		return fmt.Errorf("failed to create shuffle event: %w", err)
	}

	// Apply shuffle event
	if err := game.ApplyEvent(shuffleEvent); err != nil {
		return fmt.Errorf("failed to apply shuffle event: %w", err)
	}

	// 3. Deal cards to each player with specific card IDs
	// Get currently used cards from Redis
	usedCardIDs, err := h.EventService.GetUsedCards(h.Payload.GameID)

	if err != nil {
		return fmt.Errorf("failed to get used cards from Redis: %w", err)
	}

	// Convert to map for efficient lookup
	usedCardsMap := make(map[string]bool)

	for _, cardID := range usedCardIDs {
		usedCardsMap[cardID] = true
	}

	// Collect all cards to be used for batching Redis operations
	allUsedCardIDs := []string{}

	for _, player := range game.Players {
		// Find 7 white cards in the collection that haven't been used yet
		whiteCards := []*domain.Card{}
		for _, card := range game.Collection.Cards {
			if card.Type == "White" && len(whiteCards) < 7 && !usedCardsMap[card.ID] {
				whiteCards = append(whiteCards, card)
				usedCardsMap[card.ID] = true                     // Mark this card as used
				allUsedCardIDs = append(allUsedCardIDs, card.ID) // Collect for batch operation
			}
		}

		// Extract card IDs for the event
		cardIDs := make([]string, len(whiteCards))
		for j, card := range whiteCards {
			cardIDs[j] = card.ID
		}

		// Create deal cards event
		dealEvent, err := h.EventService.CreateGameEvent(
			h.Payload.GameID,
			domain.EventDealCards,
			domain.NewGameEventPayloadDealCards(h.Payload.GameID, player.UserID, cardIDs),
		)

		if err != nil {
			return fmt.Errorf("failed to create deal cards event: %w", err)
		}

		// Apply deal cards event
		if err := game.ApplyEvent(dealEvent); err != nil {
			return fmt.Errorf("failed to apply deal cards event: %w", err)
		}

		// Persist deal cards event
		if err := h.EventService.AppendEvent(dealEvent); err != nil {
			return fmt.Errorf("failed to persist deal cards event: %w", err)
		}
	}

	// 4. Draw black card with specific card ID
	blackCards := []*domain.Card{}
	for _, card := range game.Collection.Cards {
		if card.Type == "Black" && len(blackCards) < 1 && !usedCardsMap[card.ID] {
			blackCards = append(blackCards, card)
			usedCardsMap[card.ID] = true                     // Mark this card as used
			allUsedCardIDs = append(allUsedCardIDs, card.ID) // Collect for batch operation
		}
	}

	if len(blackCards) > 0 {
		blackCardEvent, err := h.EventService.CreateGameEvent(
			h.Payload.GameID,
			domain.EventDrawBlackCard,
			domain.NewGameEventPayloadDrawBlackCard(h.Payload.GameID, blackCards[0].ID),
		)

		if err != nil {
			return fmt.Errorf("failed to create draw black card event: %w", err)
		}

		// Apply draw black card event
		if err := game.ApplyEvent(blackCardEvent); err != nil {
			return fmt.Errorf("failed to apply draw black card event: %w", err)
		}

		// Persist draw black card event
		if err := h.EventService.AppendEvent(blackCardEvent); err != nil {
			return fmt.Errorf("failed to persist draw black card event: %w", err)
		}
	}

	// Batch add all used cards to Redis in a single operation
	if len(allUsedCardIDs) > 0 {
		if err := h.EventService.AddUsedCards(h.Payload.GameID, allUsedCardIDs); err != nil {
			return fmt.Errorf("failed to add used cards to Redis: %w", err)
		}
	}

	// Persist game begins and shuffle events
	if err := h.EventService.AppendEvent(shuffleEvent); err != nil {
		return fmt.Errorf("failed to persist shuffle event: %w", err)
	}
	if err := h.EventService.AppendEvent(gameBeginsEvent); err != nil {
		return fmt.Errorf("failed to persist game begins event: %w", err)
	}

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)
	chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, "Game has begun.")

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

	return nil
}
