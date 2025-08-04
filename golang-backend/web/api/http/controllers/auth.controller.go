package controller

import (
	"cardgame/bootstrap"
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/repositories"
	"cardgame/internal/domain/services"
	"cardgame/internal/interfaces/http/response"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
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

type AuthController struct {
	Env            *bootstrap.Env
	GoogleConfig   oauth2.Config
	DiscordConfig  oauth2.Config
	JWTAuthService services.JwtAuthService
	UserRepository repositories.UserRepository
}

func (ac *AuthController) HandleBeginGoogleOAuthLogin(c *fiber.Ctx) error {
	// Check if we're in local development and bypass is enabled
	if ac.Env.AppEnv == "development" && ac.Env.LocalDevBypass {
		log.Printf("🔧 [LOCAL_DEV] Bypassing Google OAuth for local development")

		// Create local dev user
		localDevUser := &entities.User{
			ID:            ac.Env.LocalDevUserID,
			Name:          ac.Env.LocalDevUserName,
			Email:         ac.Env.LocalDevUserEmail,
			EmailVerified: true,
			Provider:      "local-dev",
			Image:         "https://via.placeholder.com/150/4CAF50/FFFFFF?text=DEV",
		}

		// Upsert user to database
		result, err := ac.UserRepository.UpsertUserByID(localDevUser)
		if err != nil {
			log.Printf("⚠️ [LOCAL_DEV] Failed to upsert local dev user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create local dev user",
			})
		}

		// Create access token
		accessToken, err := ac.JWTAuthService.CreateAccessToken(
			result.Name,
			result.Email,
			result.ID,
			result.Image,
			result.EmailVerified,
		)
		if err != nil {
			log.Printf("❌ [LOCAL_DEV] Failed to create access token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create local dev token",
			})
		}

		// Create refresh token
		refreshToken, err := ac.JWTAuthService.CreateRefreshToken(result.ID)
		if err != nil {
			log.Printf("⚠️ [LOCAL_DEV] Failed to create refresh token: %v", err)
		}

		// Set tokens in cookies
		ac.JWTAuthService.HandleSetAccessTokenInCookie(c, accessToken)
		if refreshToken != "" {
			ac.JWTAuthService.HandleSetRefreshTokenInCookie(c, refreshToken)
		}

		log.Printf("✅ [LOCAL_DEV] Local dev user authenticated: %s (%s)", result.Name, result.Email)

		res := response.BeginAuthLoginProcess{
			RedirectURL: "http://localhost:5173/games",
		}

		// Redirect to games page
		return c.JSON(res)
	}

	// Regular OAuth flow for production/development without bypass
	state := uuid.New().String()

	googleClientID := ac.Env.GoogleOAuthClientID
	redirectURI := ac.Env.GoogleOAuthRedirectURI
	scope := "profile email"

	oauthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		googleClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	res := response.BeginAuthLoginProcess{
		RedirectURL: oauthURL,
	}

	c.Cookie(&fiber.Cookie{
		Name:     "google_oauth_state",
		Value:    state,
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})

	return c.JSON(res)
}

func (ac *AuthController) HandleGoogleAuthCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	cookieState := c.Cookies("google_oauth_state")

	if state != cookieState {
		return c.SendString("State Mismatch")
	}

	// Clear the state cookie with matching parameters
	c.Cookie(&fiber.Cookie{
		Name:     "google_oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})

	googleConfig := ac.GoogleConfig

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.SendString("Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return c.SendString("User Data Fetch Failed")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	var userData GoogleUserData

	err = json.Unmarshal(body, &userData)

	if err != nil {
		return fmt.Errorf("failed to unmarshal body into GoogleUserData: %v", err)
	}

	usr := entities.User{
		ID:            userData.ID,
		Name:          userData.GivenName,
		Email:         userData.Email,
		EmailVerified: userData.VerifiedEmail,
		Image:         userData.Picture,
		Provider:      "google",
	}

	// Block login if email exists with another provider
	existingUser, err := ac.UserRepository.GetUserByEmail(userData.Email)

	if err == nil && existingUser != nil {
		if existingUser.Provider != "google" {
			return c.Status(fiber.StatusForbidden).SendString(
				fmt.Sprintf("An account with this email already exists via %s. Please log in with that provider.", existingUser.Provider),
			)
		}
	}

	result, err := ac.UserRepository.UpsertUserByID(&usr)

	if err != nil {
		return fmt.Errorf("there was an issue getting or creating a user: %v", err)
	}

	// Create both access and refresh tokens
	accessToken, err := ac.JWTAuthService.CreateAccessToken(result.Name, result.Email, result.ID, result.Image, result.EmailVerified)
	if err != nil {
		return fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, err := ac.JWTAuthService.CreateRefreshToken(result.ID)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %v", err)
	}

	// Set both tokens in cookies
	ac.JWTAuthService.HandleSetAccessTokenInCookie(c, accessToken)
	ac.JWTAuthService.HandleSetRefreshTokenInCookie(c, refreshToken)

	return c.Redirect("http://localhost:5173/games")
}

func (ac *AuthController) HandleBeginDiscordOAuthLogin(c *fiber.Ctx) error {
	state := uuid.New().String()

	discordClientID := ac.Env.DiscordOAuthClientID
	redirectURI := ac.Env.DiscordOAuthRedirectURI
	scope := "identify email"

	oauthURL := fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		discordClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	res := response.BeginAuthLoginProcess{
		RedirectURL: oauthURL,
	}

	c.Cookie(&fiber.Cookie{
		Name:     "discord_oauth_state",
		Value:    state,
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})

	return c.JSON(res)
}

func (ac *AuthController) HandleDiscordAuthCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	cookieState := c.Cookies("discord_oauth_state")

	if state != cookieState {
		return c.SendString("State Mismatch")
	}

	// Clear the state cookie with matching parameters
	c.Cookie(&fiber.Cookie{
		Name:     "discord_oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})

	discordConfig := ac.DiscordConfig

	token, err := discordConfig.Exchange(context.Background(), code)

	if err != nil {
		return c.SendString("Code-Token Exchange Failed")
	}

	// Fetch user info from Discord
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return c.SendString("User Data Fetch Failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Discord user struct
	type DiscordUserData struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
	}

	var userData DiscordUserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body into DiscordUserData: %v", err)
	}

	// Build Discord avatar URL (optional)
	avatarURL := ""
	if userData.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userData.ID, userData.Avatar)
	}

	// Block login if email exists with another provider
	existingUser, err := ac.UserRepository.GetUserByEmail(userData.Email)
	if err == nil && existingUser != nil {
		if existingUser.Provider != "discord" {
			return c.Status(fiber.StatusForbidden).SendString(
				fmt.Sprintf("An account with this email already exists via %s. Please log in with that provider.", existingUser.Provider),
			)
		}
	}

	usr := entities.User{
		ID:            userData.ID,
		Name:          userData.Username + "#" + userData.Discriminator,
		Email:         userData.Email,
		EmailVerified: userData.Verified,
		Image:         avatarURL,
		Provider:      "discord",
	}

	result, err := ac.UserRepository.UpsertUserByID(&usr)

	if err != nil {
		return fmt.Errorf("there was an issue getting or creating a user: %v", err)
	}

	// Create both access and refresh tokens
	accessToken, err := ac.JWTAuthService.CreateAccessToken(result.Name, result.Email, result.ID, result.Image, result.EmailVerified)
	if err != nil {
		return fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, err := ac.JWTAuthService.CreateRefreshToken(result.ID)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %v", err)
	}

	// Set both tokens in cookies
	ac.JWTAuthService.HandleSetAccessTokenInCookie(c, accessToken)
	ac.JWTAuthService.HandleSetRefreshTokenInCookie(c, refreshToken)

	return c.Redirect("http://localhost:5173/games")
}

func (ac *AuthController) HandleRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("neural_decks_refresh")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No refresh token provided",
		})
	}

	newAccessToken, err := ac.JWTAuthService.RefreshAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Set new access token in cookie
	ac.JWTAuthService.HandleSetAccessTokenInCookie(c, newAccessToken)

	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
	})
}

func (ac *AuthController) HandleLogout(c *fiber.Ctx) error {
	// Get the refresh token from cookies
	refreshToken := c.Cookies("neural_decks_refresh")

	// Invalidate the refresh token if it exists
	if refreshToken != "" {
		if err := ac.JWTAuthService.InvalidateRefreshToken(refreshToken); err != nil {
			log.Printf("❌ [AUTH] Failed to invalidate refresh token during logout: %v", err)
		} else {
			log.Printf("✅ [AUTH] Refresh token invalidated during logout")
		}
	}

	// Clear all auth cookies
	ac.JWTAuthService.ClearAuthCookies(c)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (ac *AuthController) HandleInvalidateAllTokens(c *fiber.Ctx) error {
	// Get user ID from the authenticated user (from JWT claims)
	userID := c.Locals("user").(*entities.CustomClaim).UserID

	if err := ac.JWTAuthService.InvalidateAllRefreshTokensForUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to invalidate tokens",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All tokens invalidated successfully",
	})
}
