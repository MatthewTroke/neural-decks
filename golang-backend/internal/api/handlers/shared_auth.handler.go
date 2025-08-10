package handlers

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/services"
)

type SharedAuthHandler struct {
	authService *services.AuthService
}

func NewSharedAuthHandler(authService *services.AuthService) *SharedAuthHandler {
	return &SharedAuthHandler{
		authService: authService,
	}
}

func (h *SharedAuthHandler) HandleLocalDevAuth(localDevUserID, localDevUserName, localDevUserEmail, secret string) (*entities.AuthResult, error) {
	// Create local dev user
	localDevUser := &entities.User{
		ID:            localDevUserID,
		Name:          localDevUserName,
		Email:         localDevUserEmail,
		EmailVerified: true,
		Provider:      "local-dev",
		Image:         "https://via.placeholder.com/150/4CAF50/FFFFFF?text=DEV",
	}

	return h.authService.AuthenticateUser(localDevUser, secret)
}

func (h *SharedAuthHandler) HandleLogout(refreshToken string) error {
	return h.authService.LogoutUser(refreshToken)
}

func (h *SharedAuthHandler) HandleInvalidateAllTokens(userID string) error {
	return h.authService.InvalidateAllUserTokens(userID)
}
