package services

import (
	"cardgame/bootstrap/environment"
	"cardgame/domain/entities"
	"cardgame/domain/repositories"
	"fmt"
	"log"
	"time"

	"cardgame/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Ideally remove this service and add these methods as more like "util" functions
type JWTAuthService struct {
	env                    *environment.Env
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewJWTAuthService(
	env *environment.Env,
	refreshTokenRepository repositories.RefreshTokenRepository,
) *JWTAuthService {
	return &JWTAuthService{
		env:                    env,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (jwtas *JWTAuthService) CreateAccessToken(name string, email string, userId string, image string, emailVerified bool) (string, error) {
	claim := entities.CustomClaim{
		Name:          name,
		Email:         email,
		UserID:        userId,
		Image:         image,
		EmailVerified: emailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Neural Decks",
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"neural-decks-client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(jwtas.env.JWTVerifySecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jwtas *JWTAuthService) CreateRefreshToken(userId string) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    "Neural Decks",
		Subject:   userId,
		Audience:  jwt.ClaimStrings{"neural-decks-refresh"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1 hour for testing
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(jwtas.env.JWTVerifySecret))
	if err != nil {
		return "", err
	}

	hash := utils.HashToken(tokenString)
	refreshToken := &entities.RefreshToken{
		Token:     hash,
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour * 1), // 1 hour for testing
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := jwtas.refreshTokenRepository.CreateRefreshToken(refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jwtas *JWTAuthService) RefreshAccessToken(refreshToken string, user *entities.User) (string, error) {
	hash := utils.HashToken(refreshToken)
	isValid, err := jwtas.refreshTokenRepository.IsRefreshTokenValid(hash)
	if err != nil {
		return "", err
	}
	if !isValid {
		return "", fmt.Errorf("invalid or expired refresh token")
	}

	// Verify refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtas.env.JWTVerifySecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Check if it's a refresh token
	if len(claims.Audience) == 0 || claims.Audience[0] != "neural-decks-refresh" {
		return "", fmt.Errorf("not a refresh token")
	}

	// Create new access token with user details from database
	newAccessToken, err := jwtas.CreateAccessToken(user.Name, user.Email, user.ID, user.Image, user.EmailVerified)
	if err != nil {
		return "", fmt.Errorf("failed to create new access token: %v", err)
	}

	return newAccessToken, nil
}

func (jwtas *JWTAuthService) InvalidateRefreshToken(token string) error {
	hash := utils.HashToken(token)
	if err := jwtas.refreshTokenRepository.DeleteRefreshToken(hash); err != nil {
		return err
	}
	return nil
}

func (jwtas *JWTAuthService) InvalidateAllRefreshTokensForUser(userId string) error {
	if err := jwtas.refreshTokenRepository.DeleteRefreshTokensByUserID(userId); err != nil {
		return err
	}

	return nil
}

func (jwtas *JWTAuthService) VerifyJWT(tokenString string) (*entities.CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entities.CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtas.env.JWTVerifySecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*entities.CustomClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (jwtas *JWTAuthService) HandleSetAccessTokenInCookie(c *fiber.Ctx, token string) error {
	log.Printf("🍪 [JWT] Setting access token in cookie")

	// Set access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour * 7), // 7 days to match JWT expiration
		HTTPOnly: false,                              // Allow JavaScript access for frontend JWT decoding
		Secure:   false,                              // Only send over HTTPS (must be false for localhost)
		SameSite: "Lax",                              // Allow cross-site requests
		Path:     "/",                                // Available on all paths
		Domain:   "",                                 // Empty domain for localhost
	})

	return nil
}

func (jwtas *JWTAuthService) HandleSetRefreshTokenInCookie(c *fiber.Ctx, refreshToken string) error {
	log.Printf("🍪 [JWT] Setting refresh token in cookie")

	// Set refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_refresh",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * time.Hour * 30), // 30 days
		HTTPOnly: true,                                // Prevent JavaScript access for security
		Secure:   false,                               // Only send over HTTPS (must be false for localhost)
		SameSite: "Lax",                               // Allow cross-site requests
		Path:     "/",                                 // Available on all paths
		Domain:   "",                                  // Empty domain for localhost
	})

	return nil
}

func (jwtas *JWTAuthService) ClearAuthCookies(c *fiber.Ctx) error {
	log.Printf("🍪 [JWT] Clearing auth cookies")

	// Clear access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
	})

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_refresh",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
		Path:     "/",
	})

	return nil
}
