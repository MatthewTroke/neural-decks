package repositories

import (
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/domain/valueobjects"
)

type GameRepository interface {
	Create(game *aggregates.Game) error
	GetByID(id string) (*aggregates.Game, error)
	Update(game *aggregates.Game) error
	Delete(id string) error
	GetAll() ([]*aggregates.Game, error)
	GetByStatus(status valueobjects.GameStatus) ([]*aggregates.Game, error)
}
