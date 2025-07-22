package middleware

import (
	"cardgame/bootstrap"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
)

// Middleware for JWT verification
func AuthRedirect(env *bootstrap.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Cookies("neural_decks_jwt")

		if authHeader == "" {
			return c.Redirect("/login")
		}

		// TODO maybe split them up into expired / invalid as seperate things
		// Verify the JWT token
		claim, err := services.NewJWTAuthService(env).VerifyJWT(authHeader)

		if err != nil {
			return c.Redirect("/login")
		}

		// Store claims in Fiber's context (optional: to use them later in your handlers)
		c.Locals("user", claim)

		// Proceed to the next handler
		return c.Next()
	}
}
