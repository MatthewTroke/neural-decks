package auth

import "cardgame/domain/entities"

type AuthResult struct {
	User         *entities.User
	AccessToken  string
	RefreshToken string
	RedirectURL  string
}
