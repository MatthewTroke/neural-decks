package routes

import (
	"cardgame/http/controllers"
	"cardgame/http/middleware"
	"cardgame/http/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NewGameRouter(group fiber.Router, gc *controllers.GameController) {

	group.Get("/ws/game/:id",
		middleware.RequireAuth(gc.Env),
		websocket.New(gc.HandleJoinWebsocketGameRoom),
	)

	group.Get("/games", middleware.RequireAuth(gc.Env), func(c *fiber.Ctx) error {
		games, err := gc.EventService.GetAllGamesWithCurrentState()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve games with current state",
			})
		}

		response := response.GetGamesResponse(games)

		return c.Status(200).JSON(response)
	})

	group.Post("/games/new", middleware.RequireAuth(gc.Env), gc.CreateGame)
}
