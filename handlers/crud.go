package handlers

import (
	"portfolio_kurnia/database"

	"portfolio_kurnia/models"

	"github.com/gofiber/fiber/v2"
)

// CRUD generik untuk koleksi admin (pengalaman, pendidikan, testimoni, pesan, blog).
// Order dikosongkan bila tabel tidak punya kolom `order`.

// List mengembalikan semua data dalam JSON.
func List[T any](order string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []T
		q := database.DB
		if order != "" {
			q = q.Order(order)
		}
		if err := q.Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal memuat data"})
		}
		return c.JSON(items)
	}
}

// Create menyimpan data baru dari body JSON.
func Create[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item T
		if err := c.BodyParser(&item); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
		}
		if err := database.DB.Create(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data"})
		}
		return c.JSON(fiber.Map{"message": "Berhasil disimpan", "data": item})
	}
}

// Update menyimpan perubahan pada data yang sudah ada (field yang tidak
// dikirim pada JSON tetap memakai nilai lama).
func Update[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item T
		if err := database.DB.First(&item, c.Params("id")).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
		}
		if err := c.BodyParser(&item); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
		}
		if err := database.DB.Save(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan perubahan"})
		}
		return c.JSON(fiber.Map{"message": "Berhasil diperbarui", "data": item})
	}
}

// Delete menghapus data berdasarkan :id.
func Delete[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item T
		if err := database.DB.Delete(&item, c.Params("id")).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data"})
		}
		return c.JSON(fiber.Map{"message": "Berhasil dihapus"})
	}
}

// AdminStatsAPI mengembalikan jumlah objek untuk stat card & badge di dashboard.
func AdminStatsAPI(c *fiber.Ctx) error {
	var projectCount, postCount, messageCount, unreadCount int64
	database.DB.Model(&models.Project{}).Count(&projectCount)
	database.DB.Model(&models.BlogPost{}).Count(&postCount)
	database.DB.Model(&models.ContactMessage{}).Count(&messageCount)
	database.DB.Model(&models.ContactMessage{}).Where("is_read = 0").Count(&unreadCount)
	return c.JSON(fiber.Map{
		"projects": projectCount,
		"posts":    postCount,
		"messages": messageCount,
		"unread":   unreadCount,
	})
}
