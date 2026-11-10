package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

// DefaultSkills daftar keahlian yang tampil di homepage dan resume.
var DefaultSkills = []string{
	"PHP", "Laravel", "Go", "Fiber", "Python", "JavaScript",
	"Flutter", "MySQL", "PostgreSQL", "MikroTik", "Docker", "Git", "Hardware Support",
}

func Home(c *fiber.Ctx) error {
	var projects []models.Project
	var posts []models.BlogPost

	database.DB.Preload("Images").Order("`order` asc, created_at desc").Find(&projects)
	database.DB.Order("created_at desc").Limit(3).Find(&posts)

	var experience []models.Experience
	database.DB.Order("`order` asc, id asc").Find(&experience)

	var credentials []models.Credential
	database.DB.Order("`order` asc, id asc").Find(&credentials)

	var testimonials []models.Testimonial
	database.DB.Order("`order` asc, id asc").Find(&testimonials)

	settings := GetSettings()

	var menus []models.MenuLink
	database.DB.Order("`order` asc").Find(&menus)
	// Seed default menus if empty
	if len(menus) == 0 {
		defaultMenus := []models.MenuLink{
			{Label: "Beranda", URL: "#home", Order: 1},
			{Label: "Tentang", URL: "#about", Order: 2},
			{Label: "Proyek", URL: "#projects", Order: 3},
			{Label: "Kontak", URL: "#contact", Order: 4},
		}
		for _, m := range defaultMenus {
			database.DB.Create(&m)
		}
		menus = defaultMenus
	}

	data := fiber.Map{
		"Projects":      projects,
		"ProjectCount":  len(projects),
		"LatestPosts":   posts,
		"Experience":    experience,
		"Credentials":   credentials,
		"Testimonials":  testimonials,
		"ContactStatus": c.Query("kontak"),
		"Settings":      settings,
		"Menus":         menus,
		"Skills":        DefaultSkills,
	}
	data["Schema"] = HomeSchema(settings, siteURL(c))
	return c.Render("home", withSEO(c, settings, data))
}
