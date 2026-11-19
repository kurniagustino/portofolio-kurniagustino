package main

import (
	"log"
	"net/http"
	"os"
	"portfolio_kurnia/database"
	"portfolio_kurnia/handlers"
	"portfolio_kurnia/models"
	"portfolio_kurnia/templates"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	database.Connect()

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET belum diatur. Set environment variable JWT_SECRET sebelum menjalankan server.")
	}
	if os.Getenv("SITE_URL") == "" {
		log.Println("Info: SITE_URL belum diatur — canonical URL & sitemap memakai alamat request. Disarankan mengatur SITE_URL=https://domainanda.com")
	}

	// Create a new engine using embedded filesystem
	engine := html.NewFileSystem(http.FS(templates.FS), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
		// Ambil IP asli dari header proxy (Vercel/Cloudflare) untuk rate-limit.
		// Bila header tidak ada (akses lokal), fallback ke IP remote.
		ProxyHeader:        "X-Forwarded-For",
		EnableIPValidation: true,
	})

	// Hardening: recover dari panic + security headers dasar
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		ReferrerPolicy:     "strict-origin-when-cross-origin",
	}))

	// Serve static and media files
	app.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(publicFS),
		PathPrefix: "public",
		Browse:     false, // matikan directory listing
	}))
	app.Static("/media", "./media")
	app.Static("/static", "./static")

	// Public Routes
	app.Get("/", handlers.Home)
	app.Get("/projects/:slug", handlers.ProjectDetail)
	app.Get("/resume", handlers.Resume)
	app.Get("/blog", handlers.BlogIndex)
	app.Get("/blog/:slug", handlers.BlogDetail)
	app.Get("/sitemap.xml", handlers.Sitemap)
	app.Get("/robots.txt", handlers.Robots)

	// Rate-limit: maks 5 pesan per 10 menit per IP (anti spam)
	contactLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 10 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Redirect("/?kontak=terbatas#contact")
		},
	})
	app.Post("/contact", contactLimiter, handlers.ContactSubmit)

	// Auth Routes
	app.Get("/login", handlers.LoginUI)
	// Rate-limit: maks 10 percobaan login per 15 menit per IP (anti brute-force)
	loginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 15 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Terlalu banyak percobaan. Coba lagi beberapa menit lagi."})
		},
	})
	app.Post("/login", loginLimiter, handlers.LoginAPI)
	app.Get("/logout", handlers.LogoutAPI)

	// Admin API & UI Routes (Protected)
	admin := app.Group("/admin", handlers.AuthMiddleware, csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		ContextKey:     "csrf",
	}))
	admin.Get("/", handlers.AdminUI)
	admin.Get("/api/stats", handlers.AdminStatsAPI)
	admin.Get("/api/projects", handlers.AdminListProjects)
	admin.Post("/api/projects", handlers.AdminCreateProject)
	admin.Put("/api/projects/:id", handlers.AdminUpdateProject)
	admin.Delete("/api/projects/:id", handlers.AdminDeleteProject)
	admin.Delete("/api/projects/image/:id", handlers.AdminDeleteProjectImage)

	// Media upload (gambar project, blog, settings)
	admin.Post("/api/upload", handlers.UploadMediaAPI)

	// Settings API
	admin.Get("/api/settings", handlers.GetSettingsAPI)
	admin.Put("/api/settings", handlers.UpdateSettingsAPI)

	// Menus API
	admin.Get("/api/menus", handlers.GetMenusAPI)
	admin.Post("/api/menus", handlers.CreateMenuAPI)
	admin.Put("/api/menus/:id", handlers.UpdateMenuAPI)
	admin.Delete("/api/menus/:id", handlers.DeleteMenuAPI)

	// Pengalaman kerja
	admin.Get("/api/experiences", handlers.List[models.Experience]("`order` asc, id asc"))
	admin.Post("/api/experiences", handlers.Create[models.Experience]())
	admin.Put("/api/experiences/:id", handlers.Update[models.Experience]())
	admin.Delete("/api/experiences/:id", handlers.Delete[models.Experience]())

	// Pendidikan & sertifikasi
	admin.Get("/api/credentials", handlers.List[models.Credential]("`order` asc, id asc"))
	admin.Post("/api/credentials", handlers.Create[models.Credential]())
	admin.Put("/api/credentials/:id", handlers.Update[models.Credential]())
	admin.Delete("/api/credentials/:id", handlers.Delete[models.Credential]())

	// Testimoni
	admin.Get("/api/testimonials", handlers.List[models.Testimonial]("`order` asc, id asc"))
	admin.Post("/api/testimonials", handlers.Create[models.Testimonial]())
	admin.Put("/api/testimonials/:id", handlers.Update[models.Testimonial]())
	admin.Delete("/api/testimonials/:id", handlers.Delete[models.Testimonial]())

	// Blog
	admin.Get("/api/posts", handlers.List[models.BlogPost]("created_at desc"))
	admin.Post("/api/posts", handlers.AdminCreatePost)
	admin.Put("/api/posts/:id", handlers.AdminUpdatePost)
	admin.Delete("/api/posts/:id", handlers.Delete[models.BlogPost]())

	// Pesan masuk
	admin.Get("/api/messages", handlers.List[models.ContactMessage]("created_at desc"))
	admin.Put("/api/messages/:id", handlers.Update[models.ContactMessage]())
	admin.Delete("/api/messages/:id", handlers.Delete[models.ContactMessage]())

	// Determine port for Vercel or local
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Halaman 404 untuk rute yang tidak dikenal (daftarkan paling akhir)
	app.Use(func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
		}
		return c.Status(fiber.StatusNotFound).Render("not_found", fiber.Map{
			"Settings": handlers.GetSettings(),
			"Path":     c.Path(),
		})
	})

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
