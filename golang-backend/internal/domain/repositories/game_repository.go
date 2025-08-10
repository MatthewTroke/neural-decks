package repositories

import (
	"cardgame/internal/domain/aggregates"
)

type GameRepository interface {
	Create(game *aggregates.Game) (*aggregates.Game, error)
	GetAllGames() ([]*aggregates.Game, error)
	GetByID(id string) (*aggregates.Game, error)
	Update(game *aggregates.Game) (*aggregates.Game, error)
	Delete(id string) error
}
