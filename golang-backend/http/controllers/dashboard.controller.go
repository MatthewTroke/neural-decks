package controllers

import (
	"cardgame/bootstrap/environment"
	"cardgame/infra/websockets"
)

type DashboardController struct {
	Env              *environment.Env
	WebsocketService websockets.WebsocketService
}
