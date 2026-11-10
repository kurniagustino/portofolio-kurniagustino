package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

// Resume menampilkan halaman resume yang bisa dicetak / disimpan sebagai PDF.
func Resume(c *fiber.Ctx) error {
	var projects []models.Project
	database.DB.Preload("Images").Order("`order` asc, created_at desc").Find(&projects)

	var experience []models.Experience
	database.DB.Order("`order` asc, id asc").Find(&experience)

	var credentials []models.Credential
	database.DB.Order("`order` asc, id asc").Find(&credentials)

	s := GetSettings()
	return c.Render("resume", withSEO(c, s, fiber.Map{
		"Projects":    projects,
		"Experience":  experience,
		"Credentials": credentials,
		"Skills":      DefaultSkills,
		"Settings":    s,
	}))
}
