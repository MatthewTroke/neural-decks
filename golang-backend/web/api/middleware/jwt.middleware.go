package middleware

import (
	"cardgame/bootstrap"
	"cardgame/internal/infra/repositories"
	"cardgame/services"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Middleware for API routes that require authentication
func RequireAuth(env *bootstrap.Env, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the access token from cookies
		accessToken := c.Cookies("neural_decks_jwt")

		if accessToken == "" {
			log.Printf("❌ [AUTH] No access token found, returning 401")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No access token provided",
			})
		}

		// Create JWT service with user repository and refresh token repository
		jwtService := services.NewJWTAuthService(env)
		userRepo := repositories.NewSQLUserRepository(db)
		refreshTokenRepo := repositories.NewSQLRefreshTokenRepository(db)
		jwtService.SetUserRepository(userRepo)
		jwtService.SetRefreshTokenRepository(refreshTokenRepo)

		// Verify the JWT token
		claim, err := jwtService.VerifyJWT(accessToken)

		if err != nil {

			// Token is invalid, try to refresh using refresh token
			refreshToken := c.Cookies("neural_decks_refresh")

			if refreshToken != "" {
				log.Printf("🔄 [AUTH] Attempting token refresh...")
				newAccessToken, refreshErr := jwtService.RefreshAccessToken(refreshToken)
				if refreshErr == nil {
					// Set new access token
					jwtService.HandleSetAccessTokenInCookie(c, newAccessToken)

					// Verify the new token
					newClaim, newErr := jwtService.VerifyJWT(newAccessToken)
					if newErr == nil {
						log.Printf("✅ [AUTH] New token verified successfully")
						c.Locals("user", newClaim)
						return c.Next()
					}
				}
			}

			// If refresh failed or no refresh token, return unauthorized
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		log.Printf("✅ [AUTH] JWT verification successful for user: %s", claim.UserID)

		// Check if token is about to expire (within 5 minutes)
		if claim.ExpiresAt != nil {
			timeUntilExpiry := time.Until(claim.ExpiresAt.Time)

			if timeUntilExpiry < 5*time.Minute {
				log.Printf("🔄 [AUTH] Token expiring soon, attempting refresh...")
				// Token is about to expire, try to refresh
				refreshToken := c.Cookies("neural_decks_refresh")
				if refreshToken != "" {
					newAccessToken, refreshErr := jwtService.RefreshAccessToken(refreshToken)
					if refreshErr == nil {
						// Set new access token
						jwtService.HandleSetAccessTokenInCookie(c, newAccessToken)
					}
				}
			}
		}

		// Store claims in Fiber's context
		c.Locals("user", claim)
		log.Printf("✅ [AUTH] Authentication successful, proceeding to handler")

		// Proceed to the next handler
		return c.Next()
	}
}
