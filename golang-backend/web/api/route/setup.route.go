package route

import (
	"cardgame/bootstrap"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, db *gorm.DB, f *fiber.App, redis *redis.Client) {
	publicRouter := f.Group("/")
	// All Public APIs
	NewGameRouter(env, db, redis, publicRouter)
	NewAuthRouter(env, db, redis, publicRouter)
	NewDashboardRouter(env, db, redis, publicRouter)
}
