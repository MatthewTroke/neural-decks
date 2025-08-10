package entities

type AuthResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
	RedirectURL  string
}
