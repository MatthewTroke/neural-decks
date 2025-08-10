package infra

import (
	"cardgame/internal/api/controllers"
	"cardgame/internal/api/handlers"
	"cardgame/internal/infra/environment"
	"cardgame/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type DependencyContainer struct {
	Env                    *environment.Env
	DB                     *gorm.DB
	Redis                  *redis.Client
	Fiber                  *fiber.App
	GoogleConfig           oauth2.Config
	DiscordConfig          oauth2.Config
	JWTService             *services.JWTAuthService
	UserRepository         *SQLUserRepository
	RefreshTokenRepository *SQLRefreshTokenRepository
	GoogleAuthHandler      *handlers.GoogleAuthHandler
	DiscordAuthHandler     *handlers.DiscordAuthHandler
	SharedAuthHandler      *handlers.SharedAuthHandler
	AuthController         *controllers.AuthController
	GameController         *controllers.GameController
	DashboardController    *controllers.DashboardController
}

func NewDependencyContainer(
	env *environment.Env,
	db *gorm.DB,
	redis *redis.Client,
	fiber *fiber.App,
	googleConfig oauth2.Config,
	discordConfig oauth2.Config,
	jwtService *services.JWTAuthService,
	userRepository *SQLUserRepository,
	refreshTokenRepository *SQLRefreshTokenRepository,
	googleAuthHandler *handlers.GoogleAuthHandler,
	discordAuthHandler *handlers.DiscordAuthHandler,
	sharedAuthHandler *handlers.SharedAuthHandler,
	authController *controllers.AuthController,
	gameController *controllers.GameController,
	dashboardController *controllers.DashboardController,
) *DependencyContainer {

	return &DependencyContainer{
		Env:                    env,
		DB:                     db,
		Redis:                  redis,
		Fiber:                  fiber,
		GoogleConfig:           googleConfig,
		DiscordConfig:          discordConfig,
		JWTService:             jwtService,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		GoogleAuthHandler:      googleAuthHandler,
		DiscordAuthHandler:     discordAuthHandler,
		SharedAuthHandler:      sharedAuthHandler,
		AuthController:         authController,
		GameController:         gameController,
		DashboardController:    dashboardController,
	}
}
