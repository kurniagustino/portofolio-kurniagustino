package handlers

import (
	"portfolio_kurnia/database"
	"portfolio_kurnia/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "RAHASIA_NEGARA_SANGAT_AMAN" // Ganti dengan secret key yang aman (sebaiknya dari ENV)

// LoginUI merender halaman login
func LoginUI(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

// AdminUI merender halaman admin dashboard
func AdminUI(c *fiber.Ctx) error {
	return c.Render("admin", nil)
}

// Login API untuk memvalidasi user
func LoginAPI(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	var admin models.Admin
	if err := database.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau Password salah!"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau Password salah!"})
	}

	// Buat JWT Token
	claims := jwt.MapClaims{
		"id":  admin.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	// Simpan ke cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true, // Gunakan HTTPS
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Login berhasil"})
}

// Logout API
func LogoutAPI(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired ke masa lalu
		HTTPOnly: true,
	})
	return c.Redirect("/login")
}

// AuthMiddleware untuk melindungi rute admin
func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		// Jika ini request API, balas JSON. Jika web, redirect ke login.
		if c.Path() == "/admin" {
			return c.Redirect("/login")
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak. Silakan login."})
	}

	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		if c.Path() == "/admin" {
			return c.Redirect("/login")
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
	}

	return c.Next()
}
