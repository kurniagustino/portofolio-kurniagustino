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

	var settings models.SiteSetting
	database.DB.FirstOrCreate(&settings, models.SiteSetting{
		ID: 1, 
		SiteTitle: "Kurnia Gustino - Portfolio",
		HeroTitle: "IT Programmer & Software Developer",
		HeroSubtitle: "Kurnia Gustino Pratama",
		AboutText: "Terbiasa menangani pekerjaan IT secara menyeluruh...",
		FooterText: "© 2026 Kurnia Gustino",
		GithubLink: "https://github.com/kurniagustino",
		LinkedInLink: "https://www.linkedin.com/in/kurnia-gustino-pratama-17022439a/",
		EmailLink: "mailto:kurniagustino@gmail.com",
	})

	return c.Render("home", fiber.Map{
		"Projects":        projects,
		"LatestPosts":     posts,
		"Settings":        settings,
		"Skills":          []string{"PHP", "Laravel", "Go", "Fiber", "Python", "JavaScript", "Flutter", "MySQL", "PostgreSQL", "MikroTik", "Docker", "Git", "Hardware Support"},
	})
}
