package handlers

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminCreateProject(c *fiber.Ctx) error {
	// Validasi JWT sudah dilakukan di middleware
	title := c.FormValue("title")
	description := c.FormValue("description")
	techstack := c.FormValue("techstack")
	link := c.FormValue("link")
	orderStr := c.FormValue("order")
	order, _ := strconv.Atoi(orderStr)

	project := models.Project{
		Title:       title,
		Description: description,
		TechStack:   techstack,
		Link:        link,
		Order:       order,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal simpan proyek"})
	}

	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["images"]
		for _, file := range files {
			imgPath, err := saveImageWithCtx(c, file)
			if err == nil {
				img := models.ProjectImage{
					ProjectID: project.ID,
					Image:     imgPath,
				}
				database.DB.Create(&img)
			}
		}
	}

	return c.JSON(fiber.Map{"message": "Proyek berhasil disimpan"})
}

// Fungsi baru: saveImageWithCtx
func saveImageWithCtx(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	mediaDir := "./media"
	os.MkdirAll(mediaDir, 0755)
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
	path := filepath.Join(mediaDir, filename)
	if err := c.SaveFile(file, path); err != nil {
		return "", err
	}
	return filename, nil
}
