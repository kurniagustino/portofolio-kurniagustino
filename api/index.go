package api

import (
	"net/http"
	"portfolio_kurnia/database"
	"portfolio_kurnia/handlers"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var app *fiber.App

func init() {
	database.Connect()

	// Initialize the template engine
	engine := html.New("./templates", ".html")

	// Create a new Fiber app
	app = fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files
	// Note: In Vercel, static files can also be served natively by Vercel if configured,
	// but mapping them here ensures compatibility with the Fiber app.
	app.Static("/public", "./public")
	app.Static("/media", "./media")
	app.Static("/static", "./static")

	// Setup routes
	app.Get("/", handlers.Home)
	app.Get("/projects/:slug", handlers.ProjectDetail)
}

// Handler is the entry point for Vercel Serverless Functions
func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(app)(w, r)
}
