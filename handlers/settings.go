package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

// GetSettingsAPI returns the site settings
func GetSettingsAPI(c *fiber.Ctx) error {
	var setting models.SiteSetting
	if err := database.DB.Where(models.SiteSetting{ID: 1}).FirstOrCreate(&setting).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load settings"})
	}
	return c.JSON(setting)
}

// GetSettings mengambil pengaturan situs (baris tunggal ID=1)
// dengan nilai default bila belum ada.
func GetSettings() models.SiteSetting {
	var settings models.SiteSetting
	database.DB.Where(models.SiteSetting{ID: 1}).Attrs(models.SiteSetting{
		SiteTitle:    "Kurnia Gustino - Portfolio",
		HeroTitle:    "Kurnia Gustino Pratama",             // baris utama (nama)
		HeroSubtitle: "IT Programmer & Software Developer", // baris kedua (peran)
		AboutText:    "Terbiasa menangani pekerjaan IT secara menyeluruh...",
		FooterText:   "© 2026 Kurnia Gustino",
		GithubLink:   "https://github.com/kurniagustino",
		LinkedInLink: "https://www.linkedin.com/in/kurnia-gustino-pratama-17022439a/",
		EmailLink:    "mailto:kurniagustino@gmail.com",
	}).FirstOrCreate(&settings)
	return settings
}

// UpdateSettingsAPI updates the site settings
func UpdateSettingsAPI(c *fiber.Ctx) error {
	var setting models.SiteSetting
	if err := database.DB.First(&setting, 1).Error; err != nil {
		// Create if not exists
		setting.ID = 1
	}

	if err := c.BodyParser(&setting); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Force ID to 1 to maintain single row
	setting.ID = 1

	if err := database.DB.Save(&setting).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save settings"})
	}

	return c.JSON(fiber.Map{"message": "Settings updated successfully", "data": setting})
}
