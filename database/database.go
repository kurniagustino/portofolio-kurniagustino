package database

import (
	"fmt"
	"log"
	"os"
	"portfolio_kurnia/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	// Coba load .env jika ada (akan diabaikan jika tidak ada file .env, misal saat di Vercel)
	_ = godotenv.Load()

	// Kita ambil dari Environment Variable
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection successfully opened")

	// Migrate the schema (otomatis buat tabel jika belum ada)
	importedModels := []interface{}{
		&models.BlogPost{}, &models.Project{}, &models.ProjectImage{}, &models.Admin{},
		&models.SiteSetting{}, &models.MenuLink{}, &models.ContactMessage{},
		&models.Experience{}, &models.Credential{}, &models.Testimonial{},
	}
	err = DB.AutoMigrate(importedModels...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated")

	backfillProjectSlugs()
	SeedProjectsIfEmpty()
	SeedExperienceIfEmpty()
}

// backfillProjectSlugs mengisi slug kosong pada data project lama
// supaya link /projects/:slug tetap berfungsi.
func backfillProjectSlugs() {
	var projects []models.Project
	if err := DB.Where("slug = '' OR slug IS NULL").Find(&projects).Error; err != nil {
		log.Println("Skip slug backfill:", err)
		return
	}
	for _, p := range projects {
		slug := models.Slugify(p.Title)
		if slug == "untitled" {
			slug = fmt.Sprintf("project-%d", p.ID)
		}
		// Pastikan unik
		var count int64
		DB.Model(&models.Project{}).Where("slug = ? AND id <> ?", slug, p.ID).Count(&count)
		if count > 0 {
			slug = fmt.Sprintf("%s-%d", slug, p.ID)
		}
		if err := DB.Model(&models.Project{}).Where("id = ?", p.ID).Update("slug", slug).Error; err != nil {
			log.Println("Slug backfill gagal:", err)
		}
	}
	if len(projects) > 0 {
		fmt.Printf("Slug project diperbarui: %d\n", len(projects))
	}
}
