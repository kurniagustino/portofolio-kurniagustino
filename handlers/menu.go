package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

// GetMenusAPI returns all menu links
func GetMenusAPI(c *fiber.Ctx) error {
	var menus []models.MenuLink
	database.DB.Order("`order` asc").Find(&menus)
	return c.JSON(menus)
}

// CreateMenuAPI adds a new menu link
func CreateMenuAPI(c *fiber.Ctx) error {
	var menu models.MenuLink
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create menu"})
	}

	return c.JSON(fiber.Map{"message": "Menu created", "data": menu})
}

// DeleteMenuAPI deletes a menu link
func DeleteMenuAPI(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.MenuLink{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete menu"})
	}
	return c.JSON(fiber.Map{"message": "Menu deleted"})
}
