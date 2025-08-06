package repositories

import "cardgame/domain/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(id string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUser(id string) error
	UpsertUserByID(user *entities.User) (*entities.User, error)
}
