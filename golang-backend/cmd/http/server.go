package main

import (
	"cardgame/internal/api/config"
	"cardgame/internal/api/controllers"
	"cardgame/internal/api/handlers"
	"cardgame/internal/api/routes"
	"cardgame/internal/infra"
	"cardgame/internal/infra/environment"
	"cardgame/internal/infra/ws"
	"cardgame/internal/services"
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

	// WebSocket
	hub := ws.NewHub()
	go hub.Run()
	publisher := ws.NewPublisher(hub)

	// Bootstrap
	env := environment.NewEnv()

	dbArgs := infra.NewDatabaseInstanceArgs(env.DatabaseDSN)
	db := infra.NewDatabaseInstance(dbArgs)

	redisArgs := infra.NewRedisInstanceArgs(env.RedisHost, env.RedisPort, env.RedisPassword, env.RedisDB)
	redis := infra.NewRedisInstance(redisArgs)

	// Repositories
	userSqlRepository := infra.NewSQLUserRepository(db)
	refreshTokenSqlRepository := infra.NewSQLRefreshTokenRepository(db)
	gameRedisRepository := infra.NewRedisGameRepository(redis)
	eventRepository := infra.NewRedisEventRepository(redis)

	games, err := gameRedisRepository.GetAllGames()

	if err != nil {
		log.Fatalf("failed to get all games: %v", err)
	}

	// Services
	deckCreationService := services.NewChatGPTService(env.ChatGPTAPIKey)
	jwtService := services.NewJWTAuthService(refreshTokenSqlRepository)
	authService := services.NewAuthService(userSqlRepository, jwtService)

	// Game Coordinator with communication service
	gameCoordinator := services.NewGameCoordinator(
		gameRedisRepository,
		eventRepository,
		deckCreationService,
		publisher,
		games,
	)

	// Configs
	googleConfig := config.NewGoogleOAuthConfig(env.GoogleOAuthRedirectURI, env.GoogleOAuthClientID, env.GoogleOAuthClientSecret)
	discordConfig := config.NewDiscordOAuthConfig(env.DiscordOAuthRedirectURI, env.DiscordOAuthClientID, env.DiscordOAuthClientSecret)

	// Handlers
	googleAuthHandler := handlers.NewGoogleAuthHandler(authService)
	discordAuthHandler := handlers.NewDiscordAuthHandler(authService)
	sharedAuthHandler := handlers.NewSharedAuthHandler(authService)

	// Controller Container
	controllerContainer := controllers.NewControllerContainer(
		env,
		googleConfig,
		discordConfig,
		googleAuthHandler,
		discordAuthHandler,
		sharedAuthHandler,
		gameCoordinator,
		gameRedisRepository,
		hub,
	)

	routes.Setup(controllerContainer, fiberApp)

	fiberApp.Listen(":8080")
}
