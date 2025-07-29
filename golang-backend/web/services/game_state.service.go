package services

import (
	"cardgame/domain"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type GameStateService struct {
	mu           sync.RWMutex
	games        map[string]*domain.Game
	eventService *EventService
	timerTicker  *time.Ticker
	stopTimer    chan bool
	hubManager   interface {
		GetRoom(roomID string) *domain.Hub
	}
}

func NewGameStateService(eventService *EventService) *GameStateService {
	service := &GameStateService{
		games:        make(map[string]*domain.Game),
		eventService: eventService,
		stopTimer:    make(chan bool),
	}

	// Start the background timer for auto-continuation
	service.startAutoContinueTimer()

	return service
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

// SetEventService sets the EventService after creation (for circular dependency resolution)
func (s *GameStateService) SetEventService(eventService *EventService) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventService = eventService
}

// SetHubManager sets the hub manager for broadcasting messages
func (s *GameStateService) SetHubManager(hubManager interface {
	GetRoom(roomID string) *domain.Hub
}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hubManager = hubManager
}

// startAutoContinueTimer starts a background timer that checks for games that need auto-continuation
func (s *GameStateService) startAutoContinueTimer() {
	s.timerTicker = time.NewTicker(10 * time.Second) // Check every 10 seconds

	go func() {
		for {
			select {
			case <-s.timerTicker.C:
				s.checkForAutoContinue()
			case <-s.stopTimer:
				s.timerTicker.Stop()
				return
			}
		}
	}()
}

// checkForAutoContinue checks all active games and handles auto-continuation for AFK players
func (s *GameStateService) checkForAutoContinue() {
	s.mu.RLock()
	activeGames := make([]*domain.Game, 0, len(s.games))
	for _, game := range s.games {
		if game.Status == domain.InProgress && game.ShouldAutoContinue() {
			activeGames = append(activeGames, game)
		}
	}
	s.mu.RUnlock()

	for _, game := range activeGames {
		s.handleAutoContinue(game)
	}
}

// handleAutoContinue handles automatic continuation for a specific game
func (s *GameStateService) handleAutoContinue(game *domain.Game) {
	game.Lock()
	defer game.Unlock()

	// Get current game state
	currentGame, err := s.eventService.GetGameById(game.ID)
	if err != nil {
		return
	}

	switch currentGame.RoundStatus {
	case domain.PlayersPickingCard:
		s.createAutoPlayEvents(currentGame)
	case domain.CardCzarPickingWinningCard:
		s.createAutoPickEvent(currentGame)
	}
}

// createAutoPlayEvents creates EventCardPlayed events for AFK players who need to play cards
func (s *GameStateService) createAutoPlayEvents(game *domain.Game) {
	playersWhoHaventPlayed := game.GetPlayersWhoHaventPlayed()

	if len(playersWhoHaventPlayed) == 0 {
		return
	}

	// Create CardPlayed events for each AFK player
	for _, player := range playersWhoHaventPlayed {
		if len(player.Deck) == 0 {
			continue
		}

		// Pick a random card from the player's deck
		randomIndex := rand.Intn(len(player.Deck))
		randomCard := player.Deck[randomIndex]

		// Create a CardPlayed event for the auto-play
		cardPlayedEvent, err := s.eventService.CreateGameEvent(
			game.ID,
			domain.EventCardPlayed,
			domain.NewGameEventPayloadPlayCard(game.ID, randomCard.ID, &domain.CustomClaim{
				UserID: player.UserID,
				Name:   player.Name,
			}),
		)

		if err != nil {
			continue
		}

		// Apply the CardPlayed event
		if err := game.ApplyEvent(cardPlayedEvent); err != nil {
			continue
		}

		// Persist the CardPlayed event
		if err := s.eventService.AppendEvent(cardPlayedEvent); err != nil {
			continue
		}

		// Broadcast the auto-play message
		s.broadcastAutoContinueMessage(game.ID, fmt.Sprintf("🤖 %s was AFK and played a random card", player.Name))
	}
}

// createAutoPickEvent creates an EventCardCzarChoseWinningCard event for AFK card czar
func (s *GameStateService) createAutoPickEvent(game *domain.Game) {
	cardCzar, err := game.FindCurrentCardCzar()
	if err != nil {
		return
	}

	if len(game.WhiteCards) == 0 {
		return
	}

	// Pick a random white card as the winner
	randomIndex := rand.Intn(len(game.WhiteCards))
	randomCard := game.WhiteCards[randomIndex]

	// Create a CardCzarChoseWinningCard event for the auto-pick
	cardCzarChoseEvent, err := s.eventService.CreateGameEvent(
		game.ID,
		domain.EventCardCzarChoseWinningCard,
		domain.NewGameEventPayloadCardCzarChoseWinningCard(game.ID, randomCard.ID),
	)

	if err != nil {
		return
	}

	// Apply the CardCzarChoseWinningCard event
	if err := game.ApplyEvent(cardCzarChoseEvent); err != nil {
		return
	}

	// Persist the CardCzarChoseWinningCard event
	if err := s.eventService.AppendEvent(cardCzarChoseEvent); err != nil {
		return
	}

	// Broadcast the auto-pick message
	s.broadcastAutoContinueMessage(game.ID, fmt.Sprintf("🤖 %s was AFK and picked a random winning card", cardCzar.Name))
}

// broadcastAutoContinueMessage broadcasts auto-continuation messages to connected clients
func (s *GameStateService) broadcastAutoContinueMessage(gameID string, message string) {
	if s.hubManager == nil {
		return
	}

	hub := s.hubManager.GetRoom(gameID)
	if hub == nil {
		return
	}

	// Create chat message
	chatMessage := domain.NewWebSocketMessage(domain.ChatMessage, message)
	chatMessageBytes, err := json.Marshal(chatMessage)
	if err != nil {
		return
	}

	// Broadcast the message
	hub.Broadcast(chatMessageBytes)

	// Get updated game state and broadcast it
	game, err := s.eventService.GetGameById(gameID)
	if err != nil {
		return
	}

	gameUpdate := domain.NewWebSocketMessage(domain.GameUpdate, game)
	gameUpdateBytes, err := json.Marshal(gameUpdate)
	if err != nil {
		return
	}

	hub.Broadcast(gameUpdateBytes)
}

// StopTimer stops the background timer
func (s *GameStateService) StopTimer() {
	if s.timerTicker != nil {
		s.stopTimer <- true
	}
}

// TriggerAutoContinue manually triggers auto-continuation for a specific game (for testing)
func (s *GameStateService) TriggerAutoContinue(gameID string) {
	s.mu.RLock()
	game, exists := s.games[gameID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	s.handleAutoContinue(game)
}
