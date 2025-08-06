package repositories

import (
	"cardgame/domain/aggregates"
	"errors"

	"gorm.io/gorm"
)

type SQLGameRepository struct {
	db *gorm.DB
}

func NewSQLGameRepository(db *gorm.DB) *SQLGameRepository {
	return &SQLGameRepository{db: db}
}

func (r *SQLGameRepository) Create(game *aggregates.Game) error {
	return r.db.Create(game).Error
}

func (r *SQLGameRepository) GetByID(id string) (*aggregates.Game, error) {
	var game aggregates.Game
	if err := r.db.First(&game, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

func (r *SQLGameRepository) Update(game *aggregates.Game) (*aggregates.Game, error) {
	if err := r.db.Save(game).Error; err != nil {
		return nil, err
	}
	return game, nil
}

func (r *SQLGameRepository) Delete(id string) error {
	return r.db.Delete(&aggregates.Game{}, "id = ?", id).Error
}
