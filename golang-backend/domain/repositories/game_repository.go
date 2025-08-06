package repositories

import (
	"cardgame/domain/aggregates"
)

type GameRepository interface {
	Create(game *aggregates.Game) error
	GetByID(id string) (*aggregates.Game, error)
	Update(game *aggregates.Game) (*aggregates.Game, error)
	Delete(id string) error
}
