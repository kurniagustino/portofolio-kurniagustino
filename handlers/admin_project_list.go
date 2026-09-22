package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

func AdminListProjects(c *fiber.Ctx) error {
	var projects []models.Project
	database.DB.Preload("Images").Order("order asc, created_at desc").Find(&projects)
	return c.JSON(projects)
}

func AdminDeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := database.DB.Where("id = ?", id).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Proyek tidak ditemukan"})
	}
	// Hapus gambar terkait
	database.DB.Where("project_id = ?", project.ID).Delete(&models.ProjectImage{})
	// Hapus proyek
	database.DB.Delete(&project)
	return c.JSON(fiber.Map{"message": "Proyek berhasil dihapus"})
}

func AdminDeleteProjectImage(c *fiber.Ctx) error {
	id := c.Params("id")
	var img models.ProjectImage
	if err := database.DB.Where("id = ?", id).First(&img).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Gambar tidak ditemukan"})
	}
	// Hapus file fisik jika perlu
	// os.Remove("./media/" + img.Image)
	database.DB.Delete(&img)
	return c.JSON(fiber.Map{"message": "Gambar berhasil dihapus"})
}
