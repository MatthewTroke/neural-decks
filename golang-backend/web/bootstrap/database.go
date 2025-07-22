package bootstrap

import (
	"cardgame/domain"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var once sync.Once

func NewDatabaseInstance(env *Env) *gorm.DB {
	var database *gorm.DB

	dsn := env.DatabaseDSN

	if dsn == "" {
		fmt.Println("Database DSN is not set in the environment variables.")
		return nil
	}

	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			return
		}

		db.Migrator().DropTable(&domain.User{})

		err = db.AutoMigrate(
			&domain.User{},
		)

		if err != nil {
			fmt.Println("Failed to migrate tables:", err)
			return
		}

		database = db
	})

	return database
}
