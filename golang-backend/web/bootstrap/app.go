package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Env      *Env
	Postgres *gorm.DB
	Redis    *redis.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Postgres = NewDatabaseInstance(app.Env)
	app.Redis = NewRedisInstance(app.Env)
	// app.TemporalClient = DialTemporal(app.Env)

	return *app
}
