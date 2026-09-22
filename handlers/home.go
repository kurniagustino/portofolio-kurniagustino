package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	var projects []models.Project
	var posts []models.BlogPost

	database.DB.Preload("Images").Order("`order` asc, created_at desc").Find(&projects)
	database.DB.Order("created_at desc").Find(&posts)

	return c.Render("home", fiber.Map{
		"Projects":        projects,
		"LatestPosts":     posts,
		"Name":            "Kurnia Gustino Pratama",
		"Title":           "IT Programmer & Software Developer",
		"Tagline":         "IT Programmer dengan pengalaman 4+ tahun di pengembangan aplikasi, dukungan infrastruktur, dan pengelolaan data untuk sektor kesehatan, pertambangan, dan infrastruktur publik.",
		"About":           "Terbiasa menangani pekerjaan IT secara menyeluruh, mulai dari coding aplikasi (Laravel, Go, Python, Flutter), penanganan masalah teknis harian (IT support), hingga menjaga keandalan jaringan (MikroTik, Server Maintenance). Di luar urusan IT, saya juga berpengalaman langsung di lapangan dalam mengelola sewa alat berat, audit timesheet, serta memimpin survei data skala besar.",
		"Skills":          []string{"PHP", "Laravel", "Go", "Fiber", "Python", "JavaScript", "Flutter", "MySQL", "PostgreSQL", "MikroTik", "Docker", "Git", "Hardware Support"},
		"ContactEmail":    "kurniagustino@gmail.com",
		"ContactGithub":   "https://github.com/kurniagustino",
		"ContactLinkedin": "https://www.linkedin.com/in/kurnia-gustino-pratama-17022439a/",
	})
}
