package controllers

import (
	"cardgame/internal/infra/environment"
	"cardgame/internal/infra/ws"
)

type DashboardController struct {
	Env *environment.Env
	Hub *ws.Hub
}
