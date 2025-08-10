package entities

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	Name                 string `json:"name"`           // Custom field for email
	Email                string `json:"email"`          // Custom field for email
	UserID               string `json:"user_id"`        // Custom field for user ID
	Image                string `json:"image"`          // Custom field for image
	EmailVerified        bool   `json:"email_verified"` // Custom field for email verification status
	jwt.RegisteredClaims        // Embeds the standard registered claims
}
