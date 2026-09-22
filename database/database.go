package database

import (
	"fmt"
	"log"
	"portfolio_kurnia/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection successfully opened")

	// Migrate the schema (otomatis buat tabel jika belum ada)
	importedModels := []interface{}{
		&models.BlogPost{}, &models.Project{}, &models.ProjectImage{}, &models.Admin{},
	}
	err = DB.AutoMigrate(importedModels...)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated")
}
