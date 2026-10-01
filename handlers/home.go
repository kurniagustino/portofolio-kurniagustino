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
	database.DB.Where(models.SiteSetting{ID: 1}).Attrs(models.SiteSetting{
		SiteTitle: "Kurnia Gustino - Portfolio",
		HeroTitle: "IT Programmer & Software Developer",
		HeroSubtitle: "Kurnia Gustino Pratama",
		AboutText: "Terbiasa menangani pekerjaan IT secara menyeluruh...",
		FooterText: "© 2026 Kurnia Gustino",
		GithubLink: "https://github.com/kurniagustino",
		LinkedInLink: "https://www.linkedin.com/in/kurnia-gustino-pratama-17022439a/",
		EmailLink: "mailto:kurniagustino@gmail.com",
	}).FirstOrCreate(&settings)

	var menus []models.MenuLink
	database.DB.Order("`order` asc").Find(&menus)
	// Seed default menus if empty
	if len(menus) == 0 {
		defaultMenus := []models.MenuLink{
			{Label: "Home", URL: "#home", Order: 1},
			{Label: "About", URL: "#about", Order: 2},
			{Label: "Projects", URL: "#projects", Order: 3},
			{Label: "Contact", URL: "#contact", Order: 4},
		}
		for _, m := range defaultMenus {
			database.DB.Create(&m)
		}
		menus = defaultMenus
	}

	return c.Render("home", fiber.Map{
		"Projects":        projects,
		"LatestPosts":     posts,
		"Settings":        settings,
		"Menus":           menus,
		"Skills":          []string{"PHP", "Laravel", "Go", "Fiber", "Python", "JavaScript", "Flutter", "MySQL", "PostgreSQL", "MikroTik", "Docker", "Git", "Hardware Support"},
	})
}
