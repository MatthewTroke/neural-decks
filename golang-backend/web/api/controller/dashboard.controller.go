package controller

import (
	"cardgame/bootstrap"
	"cardgame/domain"
)

type DashboardController struct {
	Env              *bootstrap.Env
	GameService      domain.GameService
	WebsocketService domain.WebsocketService
}
