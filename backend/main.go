package main

import (
	"assaultrifle/Cron"
	
	"assaultrifle/Handler"
	"assaultrifle/Middleware"
	"assaultrifle/Routes"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
)

func main() {
	go Cron.StartScrapingCron()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})


	Initialize(app)

	app.Listen(":3000")
}


func Initialize(app *fiber.App) {
	app.Use(Middleware.Cors)
	app.Use(Middleware.RateLimit)
	app.Use(helmet.New())
	app.Use(Middleware.Compress)

	auth := app.Group("/api/auth")
	admin := app.Group("/api/admin")

	app.Get("/", Handler.Home)

	Routes.Auth(auth)
	Routes.Admin(admin)

	app.Use(Middleware.Notfound)
}
