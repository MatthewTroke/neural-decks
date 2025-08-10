package services

import "cardgame/internal/domain/aggregates"

type DeckCreationService interface {
	GenerateDeck(subject string) (*aggregates.Collection, error)
}
