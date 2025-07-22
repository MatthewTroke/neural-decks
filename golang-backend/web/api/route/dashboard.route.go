package route

import (
	"cardgame/bootstrap"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewDashboardRouter(env *bootstrap.Env, db *gorm.DB, group fiber.Router) {}
