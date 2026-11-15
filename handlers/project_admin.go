package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strconv"
	"strings"
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

// allowedUploadMIME: MIME yang aman untuk upload (gambar & PDF).
var allowedUploadMIME = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/gif":       ".gif",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

// maxUploadSize: 5 MB
const maxUploadSize = 5 * 1024 * 1024

// detectMIME membaca 512 byte pertama untuk deteksi tipe file (magic bytes).
func detectMIME(f *multipart.FileHeader) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	buf := make([]byte, 512)
	n, err := io.ReadAtLeast(src, buf, 1)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

// saveUploadSecure menyimpan file ke ./media dengan validasi ketat.
// - Cek ukuran (maks 5 MB)
// - Cek MIME via magic bytes, bukan ekstensi filename
// - Nama file: hex random, bukan dari input user
func saveUploadSecure(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	if file.Size > maxUploadSize {
		return "", fmt.Errorf("ukuran file melebihi batas 5 MB")
	}

	mime, err := detectMIME(file)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file")
	}
	// Potong parameter MIME, mis. "image/jpeg; charset=utf-8" → "image/jpeg"
	mime = strings.SplitN(mime, ";", 2)[0]
	mime = strings.TrimSpace(mime)

	ext, ok := allowedUploadMIME[mime]
	if !ok {
		return "", fmt.Errorf("tipe file tidak diizinkan: %s", mime)
	}

	// Validasi tambahan: pastikan PDF benar-benar diawali header %PDF
	if ext == ".pdf" {
		src, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("gagal membaca file")
		}
		header := make([]byte, 5)
		n, _ := io.ReadFull(src, header)
		src.Close()
		if n < 5 || string(header) != "%PDF-" {
			return "", fmt.Errorf("file PDF tidak valid")
		}
	}

	// Nama file: 32 byte hex random → tidak ada path traversal
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("gagal generate nama file")
	}
	filename := hex.EncodeToString(randBytes) + ext

	mediaDir := "./media"
	if err := os.MkdirAll(mediaDir, 0750); err != nil {
		return "", err
	}

	dst := filepath.Join(mediaDir, filename)
	if err := c.SaveFile(file, dst); err != nil {
		return "", err
	}
	return "/media/" + filename, nil
}

// saveImageSecure dipertahankan untuk kompatibilitas (hanya gambar).
func saveImageSecure(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	url, err := saveUploadSecure(c, file)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(url, "/media/") {
		return "", fmt.Errorf("path file tidak valid")
	}
	ext := strings.ToLower(filepath.Ext(url))
	switch ext {
	case ".jpg", ".png", ".gif", ".webp":
		return url, nil
	default:
		// hapus file non-gambar yang terlanjur tersimpan
		deleteMediaFile(url)
		return "", fmt.Errorf("hanya file gambar yang diizinkan untuk field ini")
	}
}

// saveImageWithCtx adalah alias publik (dipakai handler lain).
func saveImageWithCtx(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	return saveImageSecure(c, file)
}

// UploadMediaAPI: endpoint umum POST /admin/api/upload → { "url": "/media/xxx.ext" }
// Mendukung gambar (jpg/png/gif/webp) dan PDF (mis. untuk CV).
func UploadMediaAPI(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File tidak ditemukan di request"})
	}
	url, err := saveUploadSecure(c, file)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"url": url})
}

// deleteMediaFile menghapus file fisik dari ./media bila path adalah /media/...
func deleteMediaFile(imgPath string) {
	if !strings.HasPrefix(imgPath, "/media/") {
		return // jangan hapus file di luar ./media
	}
	// Pastikan tidak ada ".." di path
	clean := filepath.Clean(imgPath)
	if strings.Contains(clean, "..") {
		return
	}
	_ = os.Remove("." + clean)
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
