package services

import (
	"cardgame/bootstrap"
	"cardgame/domain"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthService struct {
	env *bootstrap.Env
}

func NewJWTAuthService(env *bootstrap.Env) domain.JwtAuthService {
	return &JWTAuthService{env: env}
}

func (jwtas *JWTAuthService) CreateJWT(name string, email string, userId string, image string, emailVerified bool) (string, error) {
	claim := domain.CustomClaim{
		Name:          name,
		Email:         email,
		UserID:        userId,
		Image:         image,
		EmailVerified: emailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Neural Decks",
			Subject:   "user-id-123",
			Audience:  jwt.ClaimStrings{"your-client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
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

func (jwtas *JWTAuthService) VerifyJWT(tokenString string) (*domain.CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtas.env.JWTVerifySecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.CustomClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (jwtas *JWTAuthService) HandleSetJWTInCookie(c *fiber.Ctx, token string) error {
	c.Cookie(&fiber.Cookie{
		Name:     "neural_decks_jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		HTTPOnly: false,                          // Prevent JavaScript access
		Secure:   false,                          // Only send over HTTPS (must be false for localhost)
		SameSite: "Strict",
	})

	return nil
}
