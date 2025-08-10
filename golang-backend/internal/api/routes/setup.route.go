package routes

import (
	"cardgame/internal/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(container *controllers.ControllerContainer, f *fiber.App) {
	publicRouter := f.Group("/")
	// All Public APIs
	NewGameRouter(
		publicRouter,
		container.GameController,
	)
	NewAuthRouter(
		publicRouter,
		container.AuthController,
	)
	NewDashboardRouter(
		publicRouter,
		container.DashboardController,
	)
}
