package handlers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// projectForm mengisi model Project dari input form admin.
func projectForm(c *fiber.Ctx, p *models.Project) {
	p.Title = c.FormValue("title")
	p.Description = c.FormValue("description")
	p.TechStack = c.FormValue("techstack")
	p.Link = c.FormValue("link")
	p.Year = c.FormValue("year")
	p.Category = c.FormValue("category")
	p.Role = c.FormValue("role")
	p.Problem = c.FormValue("problem")
	p.Solution = c.FormValue("solution")
	p.Result = c.FormValue("result")
	if order, err := strconv.Atoi(c.FormValue("order")); err == nil {
		p.Order = order
	}
}

func AdminCreateProject(c *fiber.Ctx) error {
	// Validasi JWT sudah dilakukan di middleware
	project := models.Project{}
	projectForm(c, &project)
	project.CreatedAt = time.Now()
	if project.Slug == "" {
		project.Slug = uniqueSlug(project.Title, 0)
	}

	if err := database.DB.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal simpan proyek"})
	}

	attachProjectImages(c, project.ID)

	return c.JSON(fiber.Map{"message": "Proyek berhasil disimpan"})
}

// AdminUpdateProject memperbarui proyek yang sudah ada (termasuk tambah gambar).
func AdminUpdateProject(c *fiber.Ctx) error {
	var project models.Project
	if err := database.DB.Preload("Images").First(&project, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Proyek tidak ditemukan"})
	}

	projectForm(c, &project)
	if project.Slug == "" {
		project.Slug = uniqueSlug(project.Title, project.ID)
	}

	if err := database.DB.Save(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal perbarui proyek"})
	}

	attachProjectImages(c, project.ID)

	return c.JSON(fiber.Map{"message": "Proyek berhasil diperbarui"})
}

// attachProjectImages menyimpan gambar yang diunggah (bila ada).
func attachProjectImages(c *fiber.Ctx, projectID uint) {
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		for _, file := range form.File["images"] {
			imgPath, err := saveImageWithCtx(c, file)
			if err == nil {
				img := models.ProjectImage{
					ProjectID: projectID,
					Image:     imgPath,
				}
				database.DB.Create(&img)
			}
		}
	}
}

// saveImageWithCtx menyimpan gambar ke ./media dan mengembalikan path publik.
func saveImageWithCtx(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	mediaDir := "./media"
	os.MkdirAll(mediaDir, 0755)
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(file.Filename)
	path := filepath.Join(mediaDir, filename)
	if err := c.SaveFile(file, path); err != nil {
		return "", err
	}
	return "/media/" + filename, nil
}

// uniqueSlug membuat slug unik dari judul, mis. "sistem-informasi-pegawai-2"
func uniqueSlug(title string, skipID uint) string {
	base := models.Slugify(title)
	slug := base
	for i := 2; ; i++ {
		q := database.DB.Model(&models.Project{}).Where("slug = ?", slug)
		if skipID > 0 {
			q = q.Where("id <> ?", skipID)
		}
		var count int64
		if err := q.Count(&count).Error; err != nil {
			return slug
		}
		if count == 0 {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}
