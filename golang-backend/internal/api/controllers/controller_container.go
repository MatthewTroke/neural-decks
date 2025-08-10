package controllers

import (
	"cardgame/internal/api/handlers"
	"cardgame/internal/domain/repositories"
	"cardgame/internal/infra/environment"
	"cardgame/internal/infra/ws"
	"cardgame/internal/services"

	"golang.org/x/oauth2"
)

type ControllerContainer struct {
	AuthController      *AuthController
	GameController      *GameController
	DashboardController *DashboardController
}

func NewControllerContainer(
	env *environment.Env,
	googleConfig oauth2.Config,
	discordConfig oauth2.Config,
	googleAuthHandler *handlers.GoogleAuthHandler,
	discordAuthHandler *handlers.DiscordAuthHandler,
	sharedAuthHandler *handlers.SharedAuthHandler,
	gameCoordinator *services.GameCoordinator,
	gameRepository repositories.GameRepository,
	hub *ws.Hub,
) *ControllerContainer {
	authController := NewAuthController(
		env,
		googleConfig,
		discordConfig,
		googleAuthHandler,
		discordAuthHandler,
		sharedAuthHandler,
	)

	gameController := NewGameController(
		env,
		gameCoordinator,
		gameRepository,
		hub,
	)

	return &ControllerContainer{
		AuthController:      authController,
		GameController:      gameController,
		DashboardController: nil,
	}
}
