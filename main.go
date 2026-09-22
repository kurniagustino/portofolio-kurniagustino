package main

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.Connect()

	// Create a new engine
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static and media files
	app.Static("/public", "./public")
	app.Static("/media", "./media")
	app.Static("/static", "./static")

	// Routes
	app.Get("/", handlers.Home)
	app.Get("/projects/:slug", handlers.ProjectDetail)

	app.Listen(":3000")
}
