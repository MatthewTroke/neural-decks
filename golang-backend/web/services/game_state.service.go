package services

import (
	"cardgame/domain"
	"fmt"
	"sync"
)

type GameStateService struct {
	mu           sync.RWMutex
	games        map[string]*domain.Game
	eventService *EventService
}

func NewGameStateService(eventService *EventService) *GameStateService {
	return &GameStateService{
		games:        make(map[string]*domain.Game),
		eventService: eventService,
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

// SetEventService sets the EventService after creation (for circular dependency resolution)
func (s *GameStateService) SetEventService(eventService *EventService) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventService = eventService
}
