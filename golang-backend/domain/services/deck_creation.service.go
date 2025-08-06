package services

import "cardgame/domain/aggregates"

type DeckCreationService interface {
	GenerateDeck(subject string) (*aggregates.Collection, error)
}
