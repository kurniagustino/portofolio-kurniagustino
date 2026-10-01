package main

import (
	"log"
	"net/http"
	"os"
	"portfolio_kurnia/database"
	"portfolio_kurnia/handlers"
	"portfolio_kurnia/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.Connect()

	// Create a new engine using embedded filesystem
	engine := html.NewFileSystem(http.FS(templates.FS), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static and media files
	// Note: in Vercel we should probably use Vercel's edge cache for these, but this works for now.
	app.Static("/public", "./public")
	app.Static("/media", "./media")
	app.Static("/static", "./static")

	// Public Routes
	app.Get("/", handlers.Home)
	app.Get("/projects/:slug", handlers.ProjectDetail)
	
	// Auth Routes
	app.Get("/login", handlers.LoginUI)
	app.Post("/login", handlers.LoginAPI)
	app.Get("/logout", handlers.LogoutAPI)

	// Admin API & UI Routes (Protected)
	admin := app.Group("/admin", handlers.AuthMiddleware)
	admin.Get("/", handlers.AdminUI)
	admin.Get("/api/projects", handlers.AdminListProjects)
	admin.Post("/api/projects", handlers.AdminCreateProject)
	admin.Delete("/api/projects/:id", handlers.AdminDeleteProject)
	admin.Delete("/api/projects/image/:id", handlers.AdminDeleteProjectImage)

	// Determine port for Vercel or local
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	app.Listen(":" + port)
}
