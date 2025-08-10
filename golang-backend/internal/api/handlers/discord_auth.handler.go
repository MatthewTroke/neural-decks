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

type DiscordUserData struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}

type DiscordAuthHandler struct {
	authService *services.AuthService
}

func NewDiscordAuthHandler(authService *services.AuthService) *DiscordAuthHandler {
	return &DiscordAuthHandler{
		authService: authService,
	}
}

func (h *DiscordAuthHandler) HandleDiscordAuth(code string, discordConfig oauth2.Config, secret string) (*entities.AuthResult, error) {
	token, err := h.authService.ExchangeOAuthCode(code, discordConfig)
	if err != nil {
		return nil, fmt.Errorf("code-token exchange failed: %w", err)
	}

	userData, err := h.fetchDiscordUserData(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Discord user data: %w", err)
	}

	avatarURL := ""
	if userData.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userData.ID, userData.Avatar)
	}

	user := &entities.User{
		ID:            userData.ID,
		Name:          userData.Username + "#" + userData.Discriminator,
		Email:         userData.Email,
		EmailVerified: userData.Verified,
		Image:         avatarURL,
		Provider:      "discord",
	}

	return h.authService.AuthenticateUser(user, secret)
}

func (h *DiscordAuthHandler) fetchDiscordUserData(accessToken string) (*DiscordUserData, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userData DiscordUserData
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &userData, nil
}
