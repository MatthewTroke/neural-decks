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
	Config         oauth2.Config
	JWTAuthService domain.JwtAuthService
	UserRepository domain.UserRepository
}

func (ac *AuthController) HandleBeginGoogleOAuthLogin(c *fiber.Ctx) error {
	state := uuid.New().String()

	googleClientID := ac.Env.GoogleOAuthClientID
	redirectURI := ac.Env.GoogleOAuthRedirectURI
	scope := "profile email"

	fmt.Println("Google OAuth Client ID:", googleClientID)
	fmt.Println("Google OAuth Redirect URI:", redirectURI)

	oauthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		googleClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(scope),
		state,
	)

	res := response.BeginGoogleAuthLoginResponse{
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

	googleConfig := ac.Config

	fmt.Println(state)
	fmt.Println(cookieState)
	fmt.Println(code)

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
