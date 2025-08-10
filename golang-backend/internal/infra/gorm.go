package infra

import (
	"cardgame/internal/domain/entities"
	"fmt"
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type DatabaseInstanceArgs struct {
	databaseDSN string
}

func NewDatabaseInstanceArgs(dsn string) DatabaseInstanceArgs {
	return DatabaseInstanceArgs{
		databaseDSN: dsn,
	}
}

func NewDatabaseInstance(args DatabaseInstanceArgs) *gorm.DB {
	var database *gorm.DB

	dsn := args.databaseDSN

	if dsn == "" {
		fmt.Println("Database DSN is not set in the environment variables.")
		return nil
	}

	log.Printf("🔍 Database Config - DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return nil
	}

	log.Println("✅ GORM Database connected successfully")

	db.Migrator().DropTable(&entities.User{})
	db.Migrator().DropTable(&entities.RefreshToken{})

	err = db.AutoMigrate(
		&entities.User{},
		&entities.RefreshToken{},
	)

	if err != nil {
		fmt.Println("Failed to migrate tables:", err)
		return nil
	}

	database = db

	return database
}
