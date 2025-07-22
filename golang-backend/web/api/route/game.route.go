package route

import (
	"cardgame/api/controller"
	"cardgame/api/middleware"
	"cardgame/bootstrap"
	"cardgame/domain"
	"cardgame/response"
	"cardgame/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

func NewGameRouter(env *bootstrap.Env, db *gorm.DB, group fiber.Router) {
	gc := &controller.GameController{
		GameService:    services.NewGameStateService(),
		RoomManager:    domain.NewRoomManager(),
		ChatGPTService: services.NewChatGPTService(env),
		Env:            env,
	}

	group.Get("/ws/game/:id",
		middleware.AuthRedirect(gc.Env),
		websocket.New(gc.HandleGameRoomWebsocket),
	)

	group.Get("/games", middleware.AuthRedirect(gc.Env), func(c *fiber.Ctx) error {
		games := gc.GameService.GetAllGames()

		response := response.GetGamesResponse(games)

		return c.Status(200).JSON(response)
	})

	group.Post("/games/new", middleware.AuthRedirect(gc.Env), gc.CreateGame)
}
