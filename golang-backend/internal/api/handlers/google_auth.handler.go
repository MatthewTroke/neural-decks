package handlers

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleUserData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GoogleAuthHandler struct {
	authService *services.AuthService
}

func NewGoogleAuthHandler(authService *services.AuthService) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		authService: authService,
	}
}

func (h *GoogleAuthHandler) HandleGoogleAuth(code string, googleConfig oauth2.Config, secret string) (*entities.AuthResult, error) {
	token, err := h.authService.ExchangeOAuthCode(code, googleConfig)
	if err != nil {
		return nil, fmt.Errorf("code-token exchange failed: %w", err)
	}

	userData, err := h.fetchGoogleUserData(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Google user data: %w", err)
	}

	// Create user entity
	user := &entities.User{
		ID:            userData.ID,
		Name:          userData.GivenName,
		Email:         userData.Email,
		EmailVerified: userData.VerifiedEmail,
		Image:         userData.Picture,
		Provider:      "google",
	}

	return h.authService.AuthenticateUser(user, secret)
}

func (h *GoogleAuthHandler) fetchGoogleUserData(accessToken string) (*GoogleUserData, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userData GoogleUserData
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &userData, nil
}
