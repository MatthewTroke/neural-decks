package route

import (
	"cardgame/api/controller"
	"cardgame/bootstrap"
	"cardgame/config"
	"cardgame/repository"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewAuthRouter(env *bootstrap.Env, db *gorm.DB, group fiber.Router) {
	sc := &controller.AuthController{
		Env:            env,
		Config:         config.NewGoogleOAuthConfig(env),
		JWTAuthService: services.NewJWTAuthService(env),
		UserRepository: repository.NewSQLUserRepository(db),
	}

	group.Get("/auth/google", sc.HandleBeginGoogleOAuthLogin)
	group.Get("/auth/google/callback", sc.HandleGoogleAuthCallback)
}
