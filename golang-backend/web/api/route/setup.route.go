package route

import (
	"cardgame/bootstrap"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, db *gorm.DB, f *fiber.App) {
	publicRouter := f.Group("/")
	// All Public APIs
	NewGameRouter(env, db, publicRouter)
	NewAuthRouter(env, db, publicRouter)
	NewDashboardRouter(env, db, publicRouter)
}
