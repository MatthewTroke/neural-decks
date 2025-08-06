package middleware

import (
	"cardgame/bootstrap/environment"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(env *environment.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("neural_decks_jwt")

		if accessToken == "" {
			log.Printf("❌ [AUTH] No access token found, returning 401")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No access token provided",
			})
		}

		token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(env.JWTVerifySecret), nil
		})

		if err != nil {
			log.Printf("❌ [AUTH] JWT parsing failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Check if token is valid
		if !token.Valid {
			log.Printf("❌ [AUTH] JWT token is invalid")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)

		if !ok {
			log.Printf("❌ [AUTH] Failed to extract claims from token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store claims in Fiber's context
		c.Locals("user", claims)
		log.Printf("✅ [AUTH] Authentication successful for user: %s, proceeding to handler", claims.Subject)

		// Proceed to the next handler
		return c.Next()
	}
}
