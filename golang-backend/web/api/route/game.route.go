package route

import (
	"cardgame/api/controller"
	"cardgame/api/middleware"
	"cardgame/bootstrap"
	"cardgame/internal/infra/repositories"
	"cardgame/internal/infra/websockets"
	"cardgame/internal/interfaces/http/response"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewGameRouter(env *bootstrap.Env, db *gorm.DB, redis *redis.Client, group fiber.Router) {
	gameStateService := services.NewGameStateService(nil)
	eventService := services.NewEventService(repositories.NewRedisEventRepository(redis), gameStateService)
	// Set EventService on GameStateService to complete the circular dependency
	gameStateService.SetEventService(eventService)

	// Create room manager for WebSocket broadcasting
	roomManager := websockets.NewRoomManager()

	// Set the room manager on the game state service for broadcasting
	gameStateService.SetRoomManager(roomManager)

	gc := controller.NewGameController(env, eventService, gameStateService, roomManager)

	group.Get("/ws/game/:id",
		middleware.RequireAuth(gc.Env, db),
		websocket.New(gc.HandleJoinWebsocketGameRoom),
	)

	group.Get("/games", middleware.RequireAuth(gc.Env, db), func(c *fiber.Ctx) error {
		games, err := gc.EventService.GetAllGamesWithCurrentState()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve games with current state",
			})
		}

		response := response.GetGamesResponse(games)

		return c.Status(200).JSON(response)
	})

	group.Post("/games/new", middleware.RequireAuth(gc.Env, db), gc.CreateGame)
}
