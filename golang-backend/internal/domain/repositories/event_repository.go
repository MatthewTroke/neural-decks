package repositories

import (
	"cardgame/internal/domain/aggregates"
	"time"
)

type EventRepository interface {
	AppendEvent(event *aggregates.GameEvent) error
	GetEventsForGame(gameID string) ([]aggregates.GameEvent, error)
	GetEventsSince(gameID string, since time.Time) ([]aggregates.GameEvent, error)
	GetEventByID(eventID string) (*aggregates.GameEvent, error)
	DeleteGameEvents(gameID string) error
	AddUsedCard(gameID, cardID string) error
	GetUsedCards(gameID string) ([]string, error)
	IsCardUsed(gameID, cardID string) (bool, error)
	ClearUsedCards(gameID string) error
	AddUsedCards(gameID string, cardIDs []string) error
}
