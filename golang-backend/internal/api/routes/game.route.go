package routes

import (
	"cardgame/internal/api/controllers"
	"cardgame/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NewGameRouter(group fiber.Router, gc *controllers.GameController) {

	group.Get("/ws/game/:id",
		middleware.RequireAuth(gc.Env),
		websocket.New(gc.HandleJoinWebsocketGameRoom),
	)

	group.Get("/games", middleware.RequireAuth(gc.Env), gc.HandleGetGames)
	group.Post("/games/new", middleware.RequireAuth(gc.Env), gc.CreateGame)
}
