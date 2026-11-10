package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ContactSubmit memproses form kontak pada homepage.
func ContactSubmit(c *fiber.Ctx) error {
	// Honeypot: robot biasanya mengisi field tersembunyi ini.
	if strings.TrimSpace(c.FormValue("website")) != "" {
		return c.Redirect("/?kontak=terkirim#contact")
	}

	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	message := strings.TrimSpace(c.FormValue("message"))

	if name == "" || message == "" {
		return c.Redirect("/?kontak=error#contact")
	}
	if len(name) > 100 || len(email) > 200 || len(message) > 5000 {
		return c.Redirect("/?kontak=error#contact")
	}

	msg := models.ContactMessage{
		Name:    name,
		Email:   email,
		Message: message,
	}
	if err := database.DB.Create(&msg).Error; err != nil {
		return c.Redirect("/?kontak=error#contact")
	}

	return c.Redirect("/?kontak=terkirim#contact")
}
