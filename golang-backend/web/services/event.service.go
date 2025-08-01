package services

import (
	"cardgame/domain"
	"cardgame/repository"
	"encoding/json"
	"fmt"
	"time"
)

type EventService struct {
	eventRepo        *repository.RedisEventRepository
	gameStateService *GameStateService
}

func NewEventService(eventRepo *repository.RedisEventRepository, gameStateService *GameStateService) *EventService {
	return &EventService{
		eventRepo:        eventRepo,
		gameStateService: gameStateService,
	}
}

func (s *EventService) AppendEvent(event domain.GameEvent) error {
	return s.eventRepo.AppendEvent(event)
}

func (s *EventService) BuildGameByGameId(gameID string) (*domain.Game, error) {
	game, err := s.gameStateService.GetGameById(gameID)

	if err != nil {
		return nil, err
	}

	return s.rebuildGameFromEvents(game)
}

// rebuildGameFromEvents rebuilds a game from all its events
func (s *EventService) rebuildGameFromEvents(initialGame *domain.Game) (*domain.Game, error) {
	events, err := s.eventRepo.GetEventsForGame(initialGame.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get events for game: %w", err)
	}

	// If no events exist, the game doesn't exist
	if len(events) == 0 {
		return initialGame, nil
	}

	// Clone game, apply all events to rebuild current state of the game.
	game := initialGame.Clone()

	for _, event := range events {
		if err := game.ApplyEvent(event); err != nil {
			return nil, fmt.Errorf("failed to apply event %s: %w", event.ID, err)
		}
	}

	return game, nil
}

func (s *EventService) CreateGameEvent(gameID string, eventType domain.GameEventType, payload interface{}) (domain.GameEvent, error) {
	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		return domain.GameEvent{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	event := domain.GameEvent{
		GameID:    gameID,
		Type:      eventType,
		Payload:   payloadJSON,
		CreatedAt: time.Now(),
	}

	return event, nil
}

func (s *EventService) GetEventsForGame(gameID string) ([]domain.GameEvent, error) {
	return s.eventRepo.GetEventsForGame(gameID)
}

func (s *EventService) GetEventsSince(gameID string, since time.Time) ([]domain.GameEvent, error) {
	return s.eventRepo.GetEventsSince(gameID, since)
}

func (s *EventService) DeleteGameEvents(gameID string) error {
	return s.eventRepo.DeleteGameEvents(gameID)
}

func (s *EventService) GetAllGamesWithCurrentState() ([]*domain.Game, error) {
	initialGames := s.gameStateService.GetAllGames()

	var currentGames []*domain.Game

	if len(initialGames) == 0 {
		return []*domain.Game{}, nil
	}

	for _, initialGame := range initialGames {
		currentGame, err := s.rebuildGameFromEvents(initialGame)
		if err != nil {
			fmt.Printf("Warning: failed to rebuild game %s from events: %v\n", initialGame.ID, err)

			currentGames = append(currentGames, initialGame)
			continue
		}
		currentGames = append(currentGames, currentGame)
	}

	return currentGames, nil
}

func (s *EventService) AddUsedCard(gameID, cardID string) error {
	return s.eventRepo.AddUsedCard(gameID, cardID)
}

func (s *EventService) AddUsedCards(gameID string, cardIDs []string) error {
	return s.eventRepo.AddUsedCards(gameID, cardIDs)
}

func (s *EventService) GetUsedCards(gameID string) ([]string, error) {
	return s.eventRepo.GetUsedCards(gameID)
}

func (s *EventService) IsCardUsed(gameID, cardID string) (bool, error) {
	return s.eventRepo.IsCardUsed(gameID, cardID)
}

func (s *EventService) ClearUsedCards(gameID string) error {
	return s.eventRepo.ClearUsedCards(gameID)
}
