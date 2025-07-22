package services

import (
	"cardgame/domain"
	"fmt"
	"sync"
)

type GameStateService struct {
	mu    sync.RWMutex
	games map[string]*domain.Game
}

func NewGameStateService() *GameStateService {
	return &GameStateService{games: make(map[string]*domain.Game)}
}

func (s *GameStateService) AddGame(game *domain.Game) *domain.Game {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
	defer s.mu.RUnlock()

	game, exists := s.games[gameId]

	if !exists {
		return nil, fmt.Errorf("could not get game by id, game with ID %s not found", gameId)
	}

	return game, nil
}
