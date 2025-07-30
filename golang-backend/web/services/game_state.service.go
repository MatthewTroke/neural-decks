package services

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameStateService struct {
	mu           sync.RWMutex
	games        map[string]*domain.Game
	eventService *EventService
	roomManager  *domain.RoomManager
	timerResets  map[string]chan bool // Channel to reset timer for each game
}

func NewGameStateService(eventService *EventService) *GameStateService {
	return &GameStateService{
		games:        make(map[string]*domain.Game),
		eventService: eventService,
		timerResets:  make(map[string]chan bool),
	}
}

func (s *GameStateService) AddGame(game *domain.Game) *domain.Game {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[game.ID] = game
	return game
}

func (s *GameStateService) GetAllGames() []*domain.Game {
	s.mu.RLock()
	defer s.mu.RUnlock()

	games := make([]*domain.Game, 0, len(s.games))
	for _, game := range s.games {
		games = append(games, game)
	}
	return games
}

func (s *GameStateService) GetGameById(gameId string) (*domain.Game, error) {
	s.mu.RLock()
	game, exists := s.games[gameId]
	s.mu.RUnlock()

	if exists {
		return game, nil
	}

	return nil, fmt.Errorf("could not get game by id, game with ID %s not found in memory", gameId)
}

func (s *GameStateService) RemoveGame(gameId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.games, gameId)
}

func (s *GameStateService) SetEventService(eventService *EventService) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventService = eventService
}

func (s *GameStateService) SetRoomManager(roomManager *domain.RoomManager) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.roomManager = roomManager
}

func (s *GameStateService) CreateAutoContinueTimer(game *domain.Game) {
	// Create reset channel for this game
	s.mu.Lock()
	s.timerResets[game.ID] = make(chan bool, 1)
	s.mu.Unlock()

	go s.autoContinueGame(game.ID)
}

// ResetAutoContinueTimer resets the auto-continue timer for a game
func (s *GameStateService) ResetAutoContinueTimer(gameID string) {
	s.mu.RLock()
	resetChan, exists := s.timerResets[gameID]
	s.mu.RUnlock()

	if exists {
		select {
		case resetChan <- true:
			log.Printf("Timer reset sent for game %s", gameID)
		default:
			// Channel is full, ignore
		}
	}
}

func (s *GameStateService) autoContinueGame(gameID string) {
	// Get reset channel for this game
	s.mu.RLock()
	resetChan := s.timerResets[gameID]
	s.mu.RUnlock()

	for {
		// Wait exactly 30 seconds or until reset
		select {
		case <-time.After(30 * time.Second):
			// 30 seconds passed, auto-progress
			game, err := s.eventService.GetGameById(gameID)
			if err != nil {
				log.Printf("Game %s not found, stopping auto-continue", gameID)
				return
			}

			if game.Status != domain.InProgress {
				log.Printf("Game %s finished, stopping auto-continue", gameID)
				return
			}

			log.Printf("Auto-progressing game %s after 30 seconds", gameID)
			s.autoProgress(game)

		case <-resetChan:
			// Timer was reset by manual action
			log.Printf("Timer reset for game %s, restarting 30-second countdown", gameID)
		}
	}
}

// autoProgress handles the automatic progression logic
func (s *GameStateService) autoProgress(game *domain.Game) {
	log.Printf("Auto-progressing game %s with round status: %s", game.ID, game.RoundStatus)

	switch game.RoundStatus {
	case domain.PlayersPickingCard:
		log.Printf("Auto-playing cards for game %s", game.ID)
		s.autoPlayCards(game)
		// Broadcast update after auto-playing
		s.broadcastGameUpdate(game, "Auto-played cards for players")
	case domain.CardCzarPickingWinningCard:
		log.Printf("Auto-picking winning card for game %s", game.ID)
		s.autoPickWinningCard(game)
		// Broadcast update after auto-picking
		s.broadcastGameUpdate(game, "Auto-picked winning card")
	case domain.CardCzarChoseWinningCard:
		log.Printf("Auto-continuing round for game %s", game.ID)
		s.autoContinueRound(game)
		// Broadcast update after auto-continuing
		s.broadcastGameUpdate(game, "Auto-continued to next round")
	default:
		log.Printf("Unknown round status: %s for game %s", game.RoundStatus, game.ID)
	}
}

// autoPlayCards automatically plays cards for players who haven't played
func (s *GameStateService) autoPlayCards(game *domain.Game) {
	for _, player := range game.Players {
		if player.IsCardCzar {
			continue
		}

		if player.PlacedCard == nil && len(player.Deck) > 0 {
			randomIndex := rand.Intn(len(player.Deck))
			cardToPlay := player.Deck[randomIndex]

			event, err := s.eventService.CreateGameEvent(
				game.ID,
				domain.EventCardPlayed,
				domain.NewGameEventPayloadPlayCard(game.ID, cardToPlay.ID, &domain.CustomClaim{UserID: player.UserID}),
			)

			if err != nil {
				log.Printf("Error creating play card event: %v", err)
				continue
			}

			if err := game.ApplyEvent(event); err != nil {
				log.Printf("Error applying play card event: %v", err)
				continue
			}

			if err := s.eventService.AppendEvent(event); err != nil {
				log.Printf("Error persisting play card event: %v", err)
				continue
			}

			log.Printf("Auto-played card for player %s", player.Name)
		}
	}
}

// autoPickWinningCard automatically picks a winning card
func (s *GameStateService) autoPickWinningCard(game *domain.Game) {
	if len(game.WhiteCards) > 0 {
		randomIndex := rand.Intn(len(game.WhiteCards))
		winningCard := game.WhiteCards[randomIndex]

		event, err := s.eventService.CreateGameEvent(
			game.ID,
			domain.EventCardCzarChoseWinningCard,
			domain.NewGameEventPayloadCardCzarChoseWinningCard(game.ID, winningCard.ID),
		)

		if err != nil {
			log.Printf("Error creating pick winning card event: %v", err)
			return
		}

		if err := game.ApplyEvent(event); err != nil {
			log.Printf("Error applying pick winning card event: %v", err)
			return
		}

		if err := s.eventService.AppendEvent(event); err != nil {
			log.Printf("Error persisting pick winning card event: %v", err)
			return
		}

		winner, err := game.FindWhiteCardOwner(winningCard)
		if err == nil {
			s.broadcastGameUpdate(game, fmt.Sprintf("Card czar chose a winning card! %s wins the round!", winner.Name))
		}
	}
}

// autoContinueRound automatically continues to the next round
func (s *GameStateService) autoContinueRound(game *domain.Game) {
	// Get used cards from Redis
	usedCardIDs, err := s.eventService.GetUsedCards(game.ID)
	if err != nil {
		log.Printf("Error getting used cards: %v", err)
		return
	}

	usedCardsMap := make(map[string]bool)
	for _, cardID := range usedCardIDs {
		usedCardsMap[cardID] = true
	}

	// Check if we need to shuffle
	playersNeedingCards := 0
	for _, player := range game.Players {
		if !player.IsCardCzar {
			playersNeedingCards++
		}
	}

	availableWhiteCards := 0
	for _, card := range game.Collection.Cards {
		if card.Type == "White" && !usedCardsMap[card.ID] {
			availableWhiteCards++
		}
	}

	// Shuffle if needed
	if availableWhiteCards < playersNeedingCards {
		if err := s.eventService.ClearUsedCards(game.ID); err != nil {
			log.Printf("Error clearing used cards: %v", err)
			return
		}

		shuffleEvent, err := s.eventService.CreateGameEvent(
			game.ID,
			domain.EventShuffle,
			domain.NewGameEventPayloadShuffle(game.ID, time.Now().UnixNano(), uuid.New().String()),
		)

		if err != nil {
			log.Printf("Error creating shuffle event: %v", err)
			return
		}

		if err := game.ApplyEvent(shuffleEvent); err != nil {
			log.Printf("Error applying shuffle event: %v", err)
			return
		}

		if err := s.eventService.AppendEvent(shuffleEvent); err != nil {
			log.Printf("Error persisting shuffle event: %v", err)
			return
		}

		s.broadcastGameUpdate(game, "No more available white cards. Re-shuffling deck...")
		usedCardsMap = make(map[string]bool)
	}

	// Determine cards for each player
	playerCards := make(map[string]string)
	newlyUsedCards := []string{}

	for _, player := range game.Players {
		if player.IsCardCzar {
			continue
		}

		for _, card := range game.Collection.Cards {
			if card.Type == "White" && !usedCardsMap[card.ID] {
				playerCards[player.UserID] = card.ID
				newlyUsedCards = append(newlyUsedCards, card.ID)
				break
			}
		}
	}

	// Check black cards
	availableBlackCards := 0
	for _, card := range game.Collection.Cards {
		if card.Type == "Black" && !usedCardsMap[card.ID] {
			availableBlackCards++
		}
	}

	if availableBlackCards < 1 {
		if err := s.eventService.ClearUsedCards(game.ID); err != nil {
			log.Printf("Error clearing used cards: %v", err)
			return
		}

		shuffleEvent, err := s.eventService.CreateGameEvent(
			game.ID,
			domain.EventShuffle,
			domain.NewGameEventPayloadShuffle(game.ID, time.Now().UnixNano(), uuid.New().String()),
		)

		if err != nil {
			log.Printf("Error creating shuffle event: %v", err)
			return
		}

		if err := game.ApplyEvent(shuffleEvent); err != nil {
			log.Printf("Error applying shuffle event: %v", err)
			return
		}

		if err := s.eventService.AppendEvent(shuffleEvent); err != nil {
			log.Printf("Error persisting shuffle event: %v", err)
			return
		}

		s.broadcastGameUpdate(game, "No more available black cards. Re-shuffling deck...")
		usedCardsMap = make(map[string]bool)
	}

	// Find black card
	var blackCardID string
	for _, card := range game.Collection.Cards {
		if card.Type == "Black" && !usedCardsMap[card.ID] {
			blackCardID = card.ID
			newlyUsedCards = append(newlyUsedCards, card.ID)
			break
		}
	}

	// Create continue round event
	event, err := s.eventService.CreateGameEvent(
		game.ID,
		domain.EventRoundContinued,
		domain.NewGameEventPayloadGameRoundContinuedWithCards(game.ID, game.Players[0].UserID, playerCards, blackCardID),
	)

	if err != nil {
		log.Printf("Error creating continue round event: %v", err)
		return
	}

	if err := game.ApplyEvent(event); err != nil {
		log.Printf("Error applying continue round event: %v", err)
		return
	}

	if len(newlyUsedCards) > 0 {
		if err := s.eventService.AddUsedCards(game.ID, newlyUsedCards); err != nil {
			log.Printf("Error adding used cards: %v", err)
			return
		}
	}

	if err := s.eventService.AppendEvent(event); err != nil {
		log.Printf("Error persisting continue round event: %v", err)
		return
	}

	s.broadcastGameUpdate(game, "Round has continued.")
}

// broadcastGameUpdate sends game updates to all connected clients
func (s *GameStateService) broadcastGameUpdate(game *domain.Game, chatMessage string) {
	if s.roomManager == nil {
		log.Printf("Warning: RoomManager not set, cannot broadcast updates")
		return
	}

	hub := s.roomManager.GetRoom(game.ID)

	message := domain.NewWebSocketMessage(domain.GameUpdate, game)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling game update: %v", err)
		return
	}

	hub.Broadcast(jsonMessage)

	if chatMessage != "" {
		chatMsg := domain.NewWebSocketMessage(domain.ChatMessage, chatMessage)
		jsonChatMessage, err := json.Marshal(chatMsg)
		if err != nil {
			log.Printf("Error marshaling chat message: %v", err)
			return
		}

		hub.Broadcast(jsonChatMessage)
	}
}
