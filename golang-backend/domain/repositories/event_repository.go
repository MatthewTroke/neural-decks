package repositories

import (
	"cardgame/domain/events"
	"time"
)

type EventRepository interface {
	AppendEvent(event *events.GameEvent) error
	GetEventsForGame(gameID string) ([]events.GameEvent, error)
	GetEventsSince(gameID string, since time.Time) ([]events.GameEvent, error)
	GetEventByID(eventID string) (*events.GameEvent, error)
	DeleteGameEvents(gameID string) error
	AddUsedCard(gameID, cardID string) error
	GetUsedCards(gameID string) ([]string, error)
	IsCardUsed(gameID, cardID string) (bool, error)
	ClearUsedCards(gameID string) error
	AddUsedCards(gameID string, cardIDs []string) error
}
