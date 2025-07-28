package main

import (
	"cardgame/api/route"
	"cardgame/bootstrap"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	err := os.MkdirAll("logs", 0755)

	if err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	file, _ := os.OpenFile("logs/logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	multiwriter := io.MultiWriter(os.Stdout, file)

	log.SetOutput(multiwriter)

	log.Println("🚀 Hot reload test - this should appear when you save!123")

	app := bootstrap.App()

	env := app.Env

	fiberApp := fiber.New()

	var CorsConfig = cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "http://localhost:5173",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           3600,
	}

	fiberApp.Use(cors.New(CorsConfig))

	route.Setup(env, app.Postgres, fiberApp, app.Redis)

	fiberApp.Listen(":8080")
}
