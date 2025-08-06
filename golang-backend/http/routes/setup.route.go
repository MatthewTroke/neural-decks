package routes

import (
	"cardgame/bootstrap/container"

	"github.com/gofiber/fiber/v2"
)

func Setup(container *container.DependencyContainer, f *fiber.App) {
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
