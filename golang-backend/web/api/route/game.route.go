package route

import (
	"cardgame/api/controller"
	"cardgame/api/middleware"
	"cardgame/bootstrap"
	"cardgame/domain"
	"cardgame/repository"
	"cardgame/response"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewGameRouter(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, group fiber.Router) {
	gameStateService := services.NewGameStateService(nil)
	eventService := services.NewEventService(repository.NewRedisEventRepository(redis), gameStateService)
	// Set EventService on GameStateService to complete the circular dependency
	gameStateService.SetEventService(eventService)

	gc := &controller.GameController{
		GameService:      gameStateService,
		EventService:     eventService,
		GameStateService: gameStateService,
		RoomManager:      domain.NewRoomManager(),
		ChatGPTService:   services.NewChatGPTService(env),
		Env:              env,
	}

	group.Get("/ws/game/:id",
		middleware.AuthRedirect(gc.Env),
		websocket.New(gc.HandleJoinWebsocketGameRoom),
	)

	group.Get("/games", middleware.AuthRedirect(gc.Env), func(c *fiber.Ctx) error {
		games, err := gc.EventService.GetAllGamesWithCurrentState()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve games with current state",
			})
		}

		response := response.GetGamesResponse(games)

		return c.Status(200).JSON(response)
	})

	group.Post("/games/new", middleware.AuthRedirect(gc.Env), gc.CreateGame)
}
