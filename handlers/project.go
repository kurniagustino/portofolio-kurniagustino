package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ProjectDetail(c *fiber.Ctx) error {
	var project models.Project
	slug := c.Params("slug")

	if err := database.DB.Preload("Images").Where("slug = ?", slug).First(&project).Error; err != nil {
		s := GetSettings()
		return c.Status(404).Render("project_detail", withSEO(c, s, fiber.Map{
			"Project":       models.Project{},
			"Settings":      s,
			"TechStackList": []string{},
			"NotFound":      true,
		}))
	}

	s := GetSettings()
	data := fiber.Map{
		"Project":       project,
		"Settings":      s,
		"TechStackList": splitTechStack(project.TechStack),
		"NotFound":      false,
	}
	// Structured data: CreativeWork + author Person
	data["Schema"] = jsonLD(map[string]interface{}{
		"@type":       "CreativeWork",
		"name":        project.Title,
		"description": project.Description,
		"url":         siteURL(c) + c.Path(),
		"inLanguage":  "id",
		"author":      personSchema(s, siteURL(c)),
		"dateCreated": project.CreatedAt.UTC().Format(time.RFC3339),
	})
	return c.Render("project_detail", withSEO(c, s, data))
}

// splitTechStack memecah "Laravel, MySQL" menjadi ["Laravel", "MySQL"]
func splitTechStack(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
