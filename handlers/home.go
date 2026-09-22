package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	var projects []models.Project
	var posts []models.BlogPost

	database.DB.Preload("Images").Order("`order` asc, created_at desc").Find(&projects)
	database.DB.Order("created_at desc").Find(&posts)

	return c.Render("home", fiber.Map{
		"Projects":        projects,
		"LatestPosts":     posts,
		"Name":            "Kurnia Gustino Pratama",
		"Title":           "Backend Developer | Laravel & Go-Fiber",
		"Tagline":         "Spesialis dalam membangun sistem backend yang scalable menggunakan Laravel dan Go-Fiber.",
		"About":           "Saya adalah seorang Backend Developer yang berfokus pada efisiensi sistem dan clean code. Keahlian utama saya terletak pada ekosistem PHP (Laravel) dan performa tinggi Go (Fiber). Selain itu, saya juga memiliki pengalaman dalam pengembangan web menggunakan Django untuk solusi yang cepat dan reliabel.",
		"Skills":          []string{"PHP", "Laravel", "Go", "Fiber", "Python", "Django", "PostgreSQL", "MySQL"},
		"ContactEmail":    "kurniagustino@gmail.com", // Sesuaikan email asli
		"ContactGithub":   "https://github.com/kurniagustino",
		"ContactLinkedin": "https://www.linkedin.com/in/kurnia-gustino-pratama-17022439a/",
	})
}
