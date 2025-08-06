package services

import (
	"cardgame/domain/aggregates"
	"cardgame/domain/entities"
	"cardgame/domain/events"
	"cardgame/domain/valueobjects"
	"cardgame/infra/websockets"
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
	games        map[string]*aggregates.Game
	eventService *EventService
	roomManager  *websockets.RoomManager
	timerResets  map[string]chan bool // Channel to reset timer for each game
}

func NewGameStateService(eventService *EventService) *GameStateService {
	return &GameStateService{
		games:        make(map[string]*aggregates.Game),
		eventService: eventService,
		timerResets:  make(map[string]chan bool),
	}
}

func (s *GameStateService) AddGame(game *aggregates.Game) *aggregates.Game {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[game.ID] = game
	return game
}

func (s *GameStateService) GetAllGames() []*aggregates.Game {
	s.mu.RLock()
	defer s.mu.RUnlock()

	games := make([]*aggregates.Game, 0, len(s.games))
	for _, game := range s.games {
		games = append(games, game)
	}
	return games
}

func (s *GameStateService) GetGameById(gameId string) (*aggregates.Game, error) {
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

func (s *GameStateService) StopGameTimer(gameId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if resetChan, exists := s.timerResets[gameId]; exists {
		select {
		case resetChan <- false: // Use false as a stop signal
		default:
			close(resetChan)
		}
		delete(s.timerResets, gameId)

		log.Printf("Stopped timer for game %s", gameId)
	}
}

// CleanupGame removes the game from memory and stops any associated timers
func (s *GameStateService) CleanupGame(gameId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove game from memory
	delete(s.games, gameId)

	// Stop timer for this game
	if resetChan, exists := s.timerResets[gameId]; exists {
		close(resetChan)
		delete(s.timerResets, gameId)
		log.Printf("Cleaned up game %s and stopped timers", gameId)
	}

	// Clean up events from Redis (do this outside the lock to avoid blocking)
	go func() {
		if err := s.eventService.DeleteGameEvents(gameId); err != nil {
			log.Printf("Failed to delete events for game %s: %v", gameId, err)
		} else {
			log.Printf("Cleaned up events for game %s", gameId)
		}
	}()
}

func (s *GameStateService) SetEventService(eventService *EventService) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventService = eventService
}

func (s *GameStateService) SetRoomManager(roomManager *websockets.RoomManager) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.roomManager = roomManager
}

func (s *GameStateService) CreateAutoContinueTimer(game *aggregates.Game) {
	s.mu.Lock()
	s.timerResets[game.ID] = make(chan bool, 1)
	s.mu.Unlock()

	go s.autoContinueGame(game.ID)
}

func (s *GameStateService) ResetAutoContinueTimer(gameID string) {
	s.mu.RLock()
	resetChan, exists := s.timerResets[gameID]
	s.mu.RUnlock()

	if exists {
		select {
		case resetChan <- true: // Send true as reset signal
			log.Printf("Timer reset sent for game %s", gameID)
		default:
			// Channel is full, ignore
		}
	}

	// Also update the clock for the game
	game, err := s.eventService.BuildGameByGameId(gameID)
	if err != nil {
		log.Printf("Error getting game for clock update: %v", err)
		return
	}

	// Update the clock to show when next auto-progress will happen
	s.updateClock(game)
}

func (s *GameStateService) autoContinueGame(gameID string) {
	// Get reset channel for this game
	s.mu.RLock()
	resetChan := s.timerResets[gameID]
	s.mu.RUnlock()

	for {
		select {
		case <-time.After(30 * time.Second):
			game, err := s.eventService.BuildGameByGameId(gameID)

			if err != nil {
				log.Printf("Game %s not found, stopping auto-continue", gameID)
				return
			}

			if game.Status != valueobjects.InProgress {
				log.Printf("Game %s finished, stopping auto-continue", gameID)
				return
			}

			log.Printf("Auto-progressing game %s after 30 seconds", gameID)
			s.autoProgress(game)

		case resetSignal := <-resetChan:
			// Check if this is a stop signal
			if !resetSignal {
				log.Printf("Timer stop signal received for game %s, stopping auto-continue", gameID)
				return
			}
			// Timer was reset by manual action
			log.Printf("Timer reset for game %s, restarting 30-second countdown", gameID)
		}
	}
}

// updateClock sets the NextAutoProgressAt timestamp to 30 seconds in the future
func (s *GameStateService) updateClock(game *aggregates.Game) {
	nextAutoProgress := time.Now().Add(30 * time.Second)

	// Create the clock update event
	event, err := s.eventService.CreateGameEvent(
		game.ID,
		events.EventClockUpdate,
		aggregates.NewGameEventPayloadClockUpdate(game.ID, nextAutoProgress),
	)

	if err != nil {
		log.Printf("Error creating clock update event: %v", err)
		return
	}

	// Apply the event to the game state
	if err := game.ApplyEvent(event); err != nil {
		log.Printf("Error applying clock update event: %v", err)
		return
	}

	// Persist the event
	if err := s.eventService.AppendEvent(event); err != nil {
		log.Printf("Error persisting clock update event: %v", err)
		return
	}

	log.Printf("Updated clock for game %s, next auto-progress at: %s", game.ID, nextAutoProgress.Format(time.RFC3339))
}

func (s *GameStateService) autoProgress(game *aggregates.Game) {
	// Update the clock to show when next auto-progress will happen
	s.updateClock(game)

	log.Printf("Auto-progressing game %s with round status: %s", game.ID, game.RoundStatus)

	switch game.RoundStatus {
	case valueobjects.PlayersPickingCard:
		log.Printf("Auto-playing cards for game %s", game.ID)
		s.autoPlayCards(game)
		s.broadcastGameUpdate(game, "Auto-played cards for players")
	case valueobjects.CardCzarPickingWinningCard:
		log.Printf("Auto-picking winning card for game %s", game.ID)
		s.autoPickWinningCard(game)
		s.broadcastGameUpdate(game, "Auto-picked winning card")
	case valueobjects.CardCzarChoseWinningCard:
		log.Printf("Auto-continuing round for game %s", game.ID)
		s.autoContinueRound(game)
		s.broadcastGameUpdate(game, "Auto-continued to next round")
	default:
		log.Printf("Unknown round status: %s for game %s", game.RoundStatus, game.ID)
	}
}

// autoPlayCards automatically plays cards for players who haven't played
func (s *GameStateService) autoPlayCards(game *aggregates.Game) {
	for _, player := range game.Players {
		if player.IsCardCzar {
			continue
		}

		if player.PlacedCard == nil && len(player.Deck) > 0 {
			randomIndex := rand.Intn(len(player.Deck))
			cardToPlay := player.Deck[randomIndex]

			event, err := s.eventService.CreateGameEvent(
				game.ID,
				events.EventCardPlayed,
				aggregates.NewGameEventPayloadPlayCard(game.ID, cardToPlay.ID, &entities.CustomClaim{UserID: player.UserID}),
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
func (s *GameStateService) autoPickWinningCard(game *aggregates.Game) {
	if len(game.WhiteCards) > 0 {
		randomIndex := rand.Intn(len(game.WhiteCards))
		winningCard := game.WhiteCards[randomIndex]

		event, err := s.eventService.CreateGameEvent(
			game.ID,
			events.EventCardCzarChoseWinningCard,
			aggregates.NewGameEventPayloadCardCzarChoseWinningCard(game.ID, winningCard.ID),
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

		// Check if anyone has won the game
		var gameWinner *aggregates.Player
		for _, player := range game.Players {
			if player.Score >= game.WinnerCount {
				gameWinner = player
				break
			}
		}

		if gameWinner != nil {
			// Game is over, someone won! Create and apply game winner event
			gameWinnerEvent, err := s.eventService.CreateGameEvent(
				game.ID,
				events.EventGameWinner,
				aggregates.NewGameEventPayloadGameWinner(game.ID, gameWinner.UserID, gameWinner.Score),
			)

			if err != nil {
				log.Printf("Error creating game winner event: %v", err)
				return
			}

			// Apply the game winner event to the game state
			if err := game.ApplyEvent(gameWinnerEvent); err != nil {
				log.Printf("Error applying game winner event: %v", err)
				return
			}

			// Persist the game winner event
			if err := s.eventService.AppendEvent(gameWinnerEvent); err != nil {
				log.Printf("Error persisting game winner event: %v", err)
				return
			}

			// Broadcast game winner message
			message := websockets.NewWebSocketMessage(aggregates.GameUpdate, game)
			chatMessage := websockets.NewWebSocketMessage(aggregates.ChatMessage, fmt.Sprintf("🎉 %s has won the game with %d points! 🎉", gameWinner.Name, gameWinner.Score))

			jsonMessage, err := json.Marshal(message)
			if err == nil {
				s.roomManager.GetRoom(game.ID).Broadcast(jsonMessage)
			}

			jsonChatMessage, err := json.Marshal(chatMessage)
			if err == nil {
				s.roomManager.GetRoom(game.ID).Broadcast(jsonChatMessage)
			}

			// Stop the timer immediately to prevent logging spam
			s.StopGameTimer(game.ID)

			// Schedule cleanup of the finished game after 30 seconds
			go func() {
				time.Sleep(30 * time.Second)
				s.CleanupGame(game.ID)
			}()
		}
	}
}

// autoContinueRound automatically continues to the next round
func (s *GameStateService) autoContinueRound(game *aggregates.Game) {
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
			events.EventShuffle,
			aggregates.NewGameEventPayloadShuffle(game.ID, time.Now().UnixNano(), uuid.New().String()),
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
			events.EventShuffle,
			aggregates.NewGameEventPayloadShuffle(game.ID, time.Now().UnixNano(), uuid.New().String()),
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
		events.EventRoundContinued,
		aggregates.NewGameEventPayloadGameRoundContinuedWithCards(game.ID, game.Players[0].UserID, playerCards, blackCardID),
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

func (s *GameStateService) broadcastGameUpdate(game *aggregates.Game, chatMessage string) {
	if s.roomManager == nil {
		log.Printf("Warning: RoomManager not set, cannot broadcast updates")
		return
	}

	hub := s.roomManager.GetRoom(game.ID)

	message := websockets.NewWebSocketMessage(aggregates.GameUpdate, game)
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling game update: %v", err)
		return
	}

	hub.Broadcast(jsonMessage)

	if chatMessage != "" {
		chatMsg := websockets.NewWebSocketMessage(aggregates.ChatMessage, chatMessage)
		jsonChatMessage, err := json.Marshal(chatMsg)
		if err != nil {
			log.Printf("Error marshaling chat message: %v", err)
			return
		}

		hub.Broadcast(jsonChatMessage)
	}
}
