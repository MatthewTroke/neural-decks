package repositories

import (
	"cardgame/domain/entities"
	"errors"

	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *SQLUserRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *SQLUserRepository) UpsertUserByID(user *entities.User) (*entities.User, error) {
	var existingUser entities.User

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

func (r *SQLUserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record
func (r *SQLUserRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

// DeleteUser removes a user from the database
func (r *SQLUserRepository) DeleteUser(id string) error {
	return r.db.Delete(&entities.User{}, "id = ?", id).Error
}
