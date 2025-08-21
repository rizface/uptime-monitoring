package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"uptime-monitoring/database"
	"uptime-monitoring/handlers"
	"uptime-monitoring/models"
	"uptime-monitoring/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	args := os.Args
	dbPath := ""

	if len(args) > 1 {
		dbPath = args[1]
	}

	// Initialize database
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize repositories and services
	urlRepo := models.NewURLRepository(db)
	checker := services.NewUptimeChecker(urlRepo)
	handler := handlers.NewHandler(urlRepo, checker)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files from embedded FS
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to create static filesystem:", err)
	}
	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}))

	// Favicon route
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		faviconData, err := staticFiles.ReadFile("static/favicon.svg")
		if err != nil {
			return c.Status(404).SendString("Favicon not found")
		}
		c.Set("Content-Type", "image/svg+xml")
		return c.Send(faviconData)
	})

	// Web routes
	app.Get("/", handler.Dashboard)
	app.Get("/add", handler.ShowAddForm)
	app.Post("/add", handler.CreateURL)
	app.Get("/edit/:id", handler.ShowEditForm)

	// API routes
	api := app.Group("/api")
	api.Get("/version", handler.Version)
	api.Get("/healthcheck", handler.Healthcheck)
	api.Get("/urls", handler.GetURLs)
	api.Get("/urls/:id", handler.GetURL)
	api.Post("/urls", handler.CreateURL)
	api.Put("/urls/:id", handler.UpdateURL)
	api.Delete("/urls/:id", handler.DeleteURL)
	api.Post("/check/:id", handler.CheckURL)
	api.Post("/check-all", handler.CheckAllURLs)

	// Handle form submissions with method override
	app.Post("/api/urls/:id", func(c *fiber.Ctx) error {
		method := c.FormValue("_method")
		if method == "PUT" {
			return handler.UpdateURL(c)
		}
		return handler.CreateURL(c)
	})

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
