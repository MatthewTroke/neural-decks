package repository

import (
	"cardgame/domain"
	"errors"

	"gorm.io/gorm"
)

type SQLGameRepository struct {
	db *gorm.DB
}

func NewSQLGameRepository(db *gorm.DB) *SQLGameRepository {
	return &SQLGameRepository{db: db}
}

func (r *SQLGameRepository) CreateGame(game *domain.Game) error {
	return r.db.Create(game).Error
}

func (r *SQLGameRepository) GetGameByID(id string) (*domain.Game, error) {
	var game domain.Game
	if err := r.db.First(&game, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

func (r *SQLGameRepository) UpdateGame(game *domain.Game) error {
	return r.db.Save(game).Error
}

func (r *SQLGameRepository) DeleteGame(id string) error {
	return r.db.Delete(&domain.Game{}, "id = ?", id).Error
}
