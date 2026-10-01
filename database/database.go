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
		&models.BlogPost{}, &models.Project{}, &models.ProjectImage{}, &models.Admin{}, &models.SiteSetting{},
	}
	err = DB.AutoMigrate(importedModels...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated")
}
