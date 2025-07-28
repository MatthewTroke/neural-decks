package controller

import (
	"cardgame/bootstrap"
	"cardgame/domain"
	"cardgame/response"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	JWTAuthService domain.JwtAuthService
	UserRepository domain.UserRepository
}

func (ac *AuthController) HandleBeginGoogleOAuthLogin(c *fiber.Ctx) error {
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

	c.ClearCookie("google_oauth_state")

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

	usr := domain.User{
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

	jwt, err := ac.JWTAuthService.CreateJWT(result.Name, result.Email, result.ID, result.Image, result.EmailVerified)

	if err != nil {
		return fmt.Errorf("failed to create jwt_: %v", err)
	}

	ac.JWTAuthService.HandleSetJWTInCookie(c, jwt)

	return c.Redirect("http://localhost:5173/dashboard")
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

	c.ClearCookie("discord_oauth_state")

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

	usr := domain.User{
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

	jwt, err := ac.JWTAuthService.CreateJWT(result.Name, result.Email, result.ID, result.Image, result.EmailVerified)

	if err != nil {
		return fmt.Errorf("failed to create jwt_: %v", err)
	}

	ac.JWTAuthService.HandleSetJWTInCookie(c, jwt)

	return c.Redirect("http://localhost:5173/dashboard")
}
