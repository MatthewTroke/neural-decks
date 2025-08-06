package main

import (
	"cardgame/bootstrap"
	"cardgame/bootstrap/container"
	"cardgame/bootstrap/environment"
	"cardgame/config"
	domainServices "cardgame/domain/services"
	"cardgame/http/controllers"
	"cardgame/http/handlers"
	"cardgame/http/routes"
	"cardgame/infra/external/ai"
	"cardgame/infra/repositories"
	"cardgame/infra/services"
	"cardgame/infra/websockets"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	err := os.MkdirAll("logs", 0755)

	if err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	file, _ := os.OpenFile("logs/logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	multiwriter := io.MultiWriter(os.Stdout, file)

	log.SetOutput(multiwriter)

	log.Println("🚀 Hot reload test - this should appear when you save!123")

	fiberApp := fiber.New()

	var CorsConfig = cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "http://localhost:5173",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           3600,
	}

	fiberApp.Use(cors.New(CorsConfig))

	// Bootstrap
	env := environment.NewEnv()
	db := bootstrap.NewDatabaseInstance(env)
	redis := bootstrap.NewRedisInstance(env)

	// Repositories
	userSqlRepository := repositories.NewSQLUserRepository(db)
	// gameSqlRepository := repositories.NewSQLGameRepository(db)
	gameRedisRepository := repositories.NewRedisGameRepository(redis)
	eventRepository := repositories.NewSQLEventRepository(db)
	refreshTokenSqlRepository := repositories.NewSQLRefreshTokenRepository(db)

	// Services
	jwtService := services.NewJWTAuthService(env, refreshTokenSqlRepository)
	authService := services.NewAuthService(userSqlRepository, jwtService)
	gameCoordinator := domainServices.NewGameCoordinator(gameRedisRepository, eventRepository)
	// eventService := services.NewEventService(redis, gameStateService)
	// gameStateService := services.NewGameStateService(eventService)
	roomManager := websockets.NewRoomManager()
	chatGPTService := ai.NewChatGPTService(env)

	// Configs
	googleConfig := config.NewGoogleOAuthConfig(env)
	discordConfig := config.NewDiscordOAuthConfig(env)

	// Handlers
	googleAuthHandler := handlers.NewGoogleAuthHandler(authService)
	discordAuthHandler := handlers.NewDiscordAuthHandler(authService)
	sharedAuthHandler := handlers.NewSharedAuthHandler(authService)

	// Controllers
	authController := controllers.NewAuthController(
		env,
		googleConfig,
		discordConfig,
		googleAuthHandler,
		discordAuthHandler,
		sharedAuthHandler,
	)
	gameController := controllers.NewGameController(
		env,
		// eventService,
		// gameStateService,
		// roomManager,
		// chatGPTService,
	)
	// dashboardController := controllers.NewDashboardController(
	// 	env,
	// )

	// Container
	container := container.NewDependencyContainer(
		env,
		db,
		redis,
		fiberApp,
		googleConfig,
		discordConfig,
		jwtService,
		userRepository,
		refreshTokenRepository,
		googleAuthHandler,
		discordAuthHandler,
		sharedAuthHandler,
		authController,
		gameController,
		dashboardController,
	)

	routes.Setup(container, fiberApp)

	fiberApp.Listen(":8080")
}
