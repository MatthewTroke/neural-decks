package route

import (
	"cardgame/bootstrap"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewDashboardRouter(env *bootstrap.Env, db *gorm.DB, redist *redis.Client, group fiber.Router) {}
