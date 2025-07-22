package bootstrap

import (
	"gorm.io/gorm"
)

type Application struct {
	Env      *Env
	Postgres *gorm.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Postgres = NewDatabaseInstance(app.Env)
	// app.TemporalClient = DialTemporal(app.Env)

	return *app
}
