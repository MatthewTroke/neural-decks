package route

import (
	"cardgame/api/middleware"
	"cardgame/bootstrap"
	"cardgame/config"
	"cardgame/internal/infra/repositories"
	controller "cardgame/internal/interfaces/http/controllers"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAuthRouter(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, group fiber.Router) {
	userRepo := repositories.NewSQLUserRepository(db)
	refreshTokenRepo := repositories.NewSQLRefreshTokenRepository(db)
	jwtService := services.NewJWTAuthService(env)
	jwtService.SetUserRepository(userRepo)
	jwtService.SetRefreshTokenRepository(refreshTokenRepo)

	sc := &controller.AuthController{
		Env:            env,
		GoogleConfig:   config.NewGoogleOAuthConfig(env),
		DiscordConfig:  config.NewDiscordOAuthConfig(env),
		JWTAuthService: jwtService,
		UserRepository: userRepo,
	}

	group.Get("/auth/google", sc.HandleBeginGoogleOAuthLogin)
	group.Get("/auth/google/callback", sc.HandleGoogleAuthCallback)
	group.Get("/auth/discord", sc.HandleBeginDiscordOAuthLogin)
	group.Get("/auth/discord/callback", sc.HandleDiscordAuthCallback)
	group.Post("/auth/refresh", sc.HandleRefreshToken)
	group.Post("/auth/logout", sc.HandleLogout)
	group.Post("/auth/invalidate-all-tokens", middleware.RequireAuth(env, db), sc.HandleInvalidateAllTokens)
}
