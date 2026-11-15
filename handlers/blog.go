package handlers

import (
	"html/template"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

var (
	ugcPolicyOnce sync.Once
	ugcPolicy     *bluemonday.Policy
)

func blogSanitizePolicy() *bluemonday.Policy {
	ugcPolicyOnce.Do(func() {
		ugcPolicy = bluemonday.UGCPolicy()
	})
	return ugcPolicy
}

var (
	tagRegex   = regexp.MustCompile(`<[^>]*>`)
	spaceRegex = regexp.MustCompile(`\s+`)
)

// BlogIndex menampilkan daftar tulisan.
func BlogIndex(c *fiber.Ctx) error {
	var posts []models.BlogPost
	database.DB.Order("created_at desc").Find(&posts)

	items := make([]fiber.Map, 0, len(posts))
	for _, p := range posts {
		items = append(items, fiber.Map{
			"Post":    p,
			"Excerpt": excerpt(p.Content, 180),
		})
	}

	s := GetSettings()
	return c.Render("blog", withSEO(c, s, fiber.Map{
		"Posts":    items,
		"Settings": s,
	}))
}

// BlogDetail menampilkan satu tulisan lengkap.
func BlogDetail(c *fiber.Ctx) error {
	var post models.BlogPost
	if err := database.DB.Where("slug = ?", c.Params("slug")).First(&post).Error; err != nil {
		s := GetSettings()
		return c.Status(404).Render("blog_detail", withSEO(c, s, fiber.Map{
			"NotFound": true,
			"Settings": s,
		}))
	}

	s := GetSettings()
	desc := excerpt(post.Content, 160)
	safeContent := blogSanitizePolicy().Sanitize(post.Content)
	data := fiber.Map{
		"Post":     post,
		"Content":  template.HTML(safeContent),
		"Excerpt":  desc,
		"Settings": s,
		"NotFound": false,
	}

	// Structured data: BlogPosting
	schema := map[string]interface{}{
		"@type":            "BlogPosting",
		"headline":         post.Title,
		"description":      desc,
		"datePublished":    post.CreatedAt.UTC().Format(time.RFC3339),
		"dateModified":     post.UpdatedAt.UTC().Format(time.RFC3339),
		"inLanguage":       "id",
		"mainEntityOfPage": siteURL(c) + c.Path(),
		"author":           personSchema(s, siteURL(c)),
		"publisher":        personSchema(s, siteURL(c)),
	}
	if post.Image != "" {
		schema["image"] = absURL(c, post.Image)
	}
	data["Schema"] = jsonLD(schema)

	return c.Render("blog_detail", withSEO(c, s, data))
}

// excerpt membuang tag HTML lalu memotong teks sampai panjang tertentu.
func excerpt(html string, max int) string {
	text := tagRegex.ReplaceAllString(html, " ")
	text = spaceRegex.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	if len(text) <= max {
		return text
	}
	cut := text[:max]
	if idx := strings.LastIndex(cut, " "); idx > 0 {
		cut = cut[:idx]
	}
	return cut + "…"
}

// AdminCreatePost membuat post blog baru dengan sanitasi konten HTML.
func AdminCreatePost(c *fiber.Ctx) error {
	var post models.BlogPost
	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	post.Content = blogSanitizePolicy().Sanitize(post.Content)
	if err := database.DB.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}
	return c.JSON(fiber.Map{"message": "Berhasil disimpan", "data": post})
}

// AdminUpdatePost memperbarui post blog dengan sanitasi konten HTML.
func AdminUpdatePost(c *fiber.Ctx) error {
	var post models.BlogPost
	if err := database.DB.First(&post, c.Params("id")).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}
	post.Content = blogSanitizePolicy().Sanitize(post.Content)
	if err := database.DB.Save(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan perubahan"})
	}
	return c.JSON(fiber.Map{"message": "Berhasil diperbarui", "data": post})
}
