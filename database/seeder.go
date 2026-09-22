package database

import (
	"fmt"
	"portfolio_kurnia/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	username := "admin"
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	admin := models.Admin{
		Username: username,
		Password: string(hash),
	}
	result := DB.Where(models.Admin{Username: username}).FirstOrCreate(&admin)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println("Seeder admin selesai. Username: admin, Password: admin123")
}
