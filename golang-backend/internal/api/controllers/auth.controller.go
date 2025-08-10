package controllers

import (
	"cardgame/internal/api/handlers"
	"cardgame/internal/api/response"
	"cardgame/internal/domain/entities"
	"cardgame/internal/infra/environment"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type AuthController struct {
	Env *environment.Env
	// TODO: Specifically make this of type googleconfig/discordconfig
	GoogleConfig       oauth2.Config
	DiscordConfig      oauth2.Config
	GoogleAuthHandler  *handlers.GoogleAuthHandler
	DiscordAuthHandler *handlers.DiscordAuthHandler
	SharedAuthHandler  *handlers.SharedAuthHandler
}

func NewAuthController(
	env *environment.Env,
	googleConfig oauth2.Config,
	discordConfig oauth2.Config,
	googleAuthHandler *handlers.GoogleAuthHandler,
	discordAuthHandler *handlers.DiscordAuthHandler,
	sharedAuthHandler *handlers.SharedAuthHandler,
) *AuthController {
	return &AuthController{
		Env:                env,
		GoogleConfig:       googleConfig,
		DiscordConfig:      discordConfig,
		GoogleAuthHandler:  googleAuthHandler,
		DiscordAuthHandler: discordAuthHandler,
		SharedAuthHandler:  sharedAuthHandler,
	}
}

func (ac *AuthController) HandleLocalDevelopmentAuth(c *fiber.Ctx) error {
	log.Printf("🔧 [LOCAL_DEV] Bypassing Auth for local development")

	// Generate random user credentials using gofakeit
	userID := fmt.Sprintf("usr_%s", gofakeit.UUID())
	userName := gofakeit.Name() + " (LOCAL_DEV)"
	userEmail := gofakeit.Email()

	log.Printf("🎲 [LOCAL_DEV] Generated random user: %s (%s) - %s", userName, userEmail, userID)

	// Delegate to shared auth handler
	result, err := ac.SharedAuthHandler.HandleLocalDevAuth(
		userID,
		userName,
		userEmail,
		ac.Env.JWTVerifySecret,
	)

	if err != nil {
		log.Printf("⚠️ [LOCAL_DEV] Failed to handle local dev auth: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create local dev user",
		})
	}

	// Set tokens in cookies
	ac.setAuthCookies(c, result.AccessToken, result.RefreshToken)

	log.Printf("✅ [LOCAL_DEV] Local dev user authenticated: %s (%s)", result.User.Name, result.User.Email)

	res := response.BeginAuthLoginProcess{
		RedirectURL: result.RedirectURL,
	}

	return c.JSON(res)
}

func (ac *AuthController) HandleBeginGoogleOAuthLogin(c *fiber.Ctx) error {
	// Check if we're in local development and bypass is enabled
	if ac.Env.AppEnv == "development" && ac.Env.LocalDevBypass {
		return ac.HandleLocalDevelopmentAuth(c)
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

	ac.setOAuthStateCookie(c, "google_oauth_state", state)

	return c.JSON(res)
}

func (ac *AuthController) HandleGoogleAuthCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")
	cookieState := c.Cookies("google_oauth_state")

	if state != cookieState {
		return c.SendString("State Mismatch")
	}

	// Clear the state cookie
	ac.clearOAuthStateCookie(c, "google_oauth_state")

	// Delegate to Google auth handler
	result, err := ac.GoogleAuthHandler.HandleGoogleAuth(code, ac.GoogleConfig, ac.Env.JWTVerifySecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Set tokens in cookies
	ac.setAuthCookies(c, result.AccessToken, result.RefreshToken)

	return c.Redirect(result.RedirectURL)
}

func (ac *AuthController) HandleBeginDiscordOAuthLogin(c *fiber.Ctx) error {
	if ac.Env.AppEnv == "development" && ac.Env.LocalDevBypass {
		return ac.HandleLocalDevelopmentAuth(c)
	}

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

	ac.setOAuthStateCookie(c, "discord_oauth_state", state)

	return c.JSON(res)
}

func (ac *AuthController) HandleDiscordAuthCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")
	cookieState := c.Cookies("discord_oauth_state")

	if state != cookieState {
		return c.SendString("State Mismatch")
	}

	// Clear the state cookie
	ac.clearOAuthStateCookie(c, "discord_oauth_state")

	// Delegate to Discord auth handler
	result, err := ac.DiscordAuthHandler.HandleDiscordAuth(code, ac.DiscordConfig, ac.Env.JWTVerifySecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Set tokens in cookies
	ac.setAuthCookies(c, result.AccessToken, result.RefreshToken)

	return c.Redirect(result.RedirectURL)
}

func (ac *AuthController) HandleRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("neural_decks_refresh_token")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No refresh token provided",
		})
	}

	// TODO: Implement refresh token logic in auth handler
	return c.JSON(fiber.Map{
		"message": "Token refreshed successfully",
	})
}

func (ac *AuthController) HandleLogout(c *fiber.Ctx) error {
	// Get the refresh token from cookies
	refreshToken := c.Cookies("neural_decks_refresh_token")

	// Delegate to shared auth handler
	if err := ac.SharedAuthHandler.HandleLogout(refreshToken); err != nil {
		log.Printf("❌ [AUTH] Failed to handle logout: %v", err)
	} else {
		log.Printf("✅ [AUTH] User logged out successfully")
	}

	// Clear all auth cookies
	ac.clearAuthCookies(c)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (ac *AuthController) HandleInvalidateAllTokens(c *fiber.Ctx) error {
	// Get user ID from the authenticated user (from JWT claims)
	userID := c.Locals("user").(*entities.CustomClaim).UserID

	// Delegate to shared auth handler
	if err := ac.SharedAuthHandler.HandleInvalidateAllTokens(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to invalidate tokens",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All tokens invalidated successfully",
	})
}

// Helper methods for HTTP concerns

func (ac *AuthController) setOAuthStateCookie(c *fiber.Ctx, name, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})
}

func (ac *AuthController) clearOAuthStateCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})
}

func (ac *AuthController) setAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	// Set access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_access_token",
		Value:    accessToken,
		HTTPOnly: false,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
		MaxAge:   3600, // 1 hour
	})

	// Set refresh token cookie
	if refreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "neural_decks_refresh_token",
			Value:    refreshToken,
			HTTPOnly: true,
			SameSite: "lax",
			Secure:   false,
			Path:     "/",
			MaxAge:   2592000, // 30 days
		})
	}
}

func (ac *AuthController) clearAuthCookies(c *fiber.Ctx) {
	// Clear access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   false,
		Path:     "/",
	})
}
