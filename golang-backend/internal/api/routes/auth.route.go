package routes

import (
	"cardgame/internal/api/controllers"
	"cardgame/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewAuthRouter(
	group fiber.Router,
	gc *controllers.AuthController,
) {
	group.Get("/auth/google", gc.HandleBeginGoogleOAuthLogin)
	group.Get("/auth/google/callback", gc.HandleGoogleAuthCallback)
	group.Get("/auth/discord", gc.HandleBeginDiscordOAuthLogin)
	group.Get("/auth/discord/callback", gc.HandleDiscordAuthCallback)
	group.Post("/auth/refresh", gc.HandleRefreshToken)
	group.Post("/auth/logout", gc.HandleLogout)
	group.Post("/auth/invalidate-all-tokens", middleware.RequireAuth(gc.Env), gc.HandleInvalidateAllTokens)
}
