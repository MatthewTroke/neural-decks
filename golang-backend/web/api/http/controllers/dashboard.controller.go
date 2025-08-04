package controller

import (
	"cardgame/bootstrap"
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/infra/websockets"
)

type DashboardController struct {
	Env              *bootstrap.Env
	GameService      aggregates.GameService
	WebsocketService websockets.WebsocketService
}
