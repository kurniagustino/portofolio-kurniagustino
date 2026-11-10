package handlers

import (
	"html/template"
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	data := fiber.Map{
		"Post":     post,
		"Content":  template.HTML(post.Content),
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
