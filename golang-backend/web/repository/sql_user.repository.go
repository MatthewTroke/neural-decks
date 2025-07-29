package repository

import (
	"cardgame/domain"
	"errors"

	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *SQLUserRepository) GetUserByID(id string) (*domain.User, error) {
	var user domain.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *SQLUserRepository) UpsertUserByID(user *domain.User) (*domain.User, error) {
	var existingUser domain.User

	// Try to find the user by ID
	err := r.db.First(&existingUser, "id = ?", user.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If not found, create a new user
			if err := r.db.Create(user).Error; err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}

	// If found, update the user with new data
	if err := r.db.Model(&existingUser).Updates(user).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}

func (r *SQLUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record
func (r *SQLUserRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

// DeleteUser removes a user from the database
func (r *SQLUserRepository) DeleteUser(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}
