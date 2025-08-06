package routes

import (
	"cardgame/http/controllers"
	"cardgame/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewAuthRouter(
	group fiber.Router,
	gc *controllers.AuthController,
) {
	// userRepo := repositories.NewSQLUserRepository(db)
	// refreshTokenRepo := repositories.NewSQLRefreshTokenRepository(db)
	// jwtService := services.NewJWTAuthService(env)
	// jwtService.SetRefreshTokenRepository(refreshTokenRepo)

	// // Create auth service
	// authService := services.NewAuthService(userRepo, jwtService)

	// // Create handlers
	// googleAuthHandler := handlers.NewGoogleAuthHandler(authService)
	// discordAuthHandler := handlers.NewDiscordAuthHandler(authService)
	// sharedAuthHandler := handlers.NewSharedAuthHandler(authService)

	// Create controller with new handlers

	group.Get("/auth/google", gc.HandleBeginGoogleOAuthLogin)
	group.Get("/auth/google/callback", gc.HandleGoogleAuthCallback)
	group.Get("/auth/discord", gc.HandleBeginDiscordOAuthLogin)
	group.Get("/auth/discord/callback", gc.HandleDiscordAuthCallback)
	group.Post("/auth/refresh", gc.HandleRefreshToken)
	group.Post("/auth/logout", gc.HandleLogout)
	group.Post("/auth/invalidate-all-tokens", middleware.RequireAuth(gc.Env), gc.HandleInvalidateAllTokens)
}
