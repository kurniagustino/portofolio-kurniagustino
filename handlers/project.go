package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

func ProjectDetail(c *fiber.Ctx) error {
	var project models.Project
	slug := c.Params("slug")

	if err := database.DB.Preload("Images").Where("slug = ?", slug).First(&project).Error; err != nil {
		return c.Status(404).SendString("Project not found")
	}

	return c.Render("project_detail", project)
}
