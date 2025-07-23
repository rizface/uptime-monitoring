package handlers

import (
	"strconv"
	"uptime-monitoring/models"
	"uptime-monitoring/services"
	"uptime-monitoring/templates"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	urlRepo *models.URLRepository
	checker *services.UptimeChecker
}

func NewHandler(urlRepo *models.URLRepository, checker *services.UptimeChecker) *Handler {
	return &Handler{
		urlRepo: urlRepo,
		checker: checker,
	}
}

// Web Routes
func (h *Handler) Dashboard(c *fiber.Ctx) error {
	urls, err := h.urlRepo.GetAll()
	if err != nil {
		return c.Status(500).SendString("Error fetching URLs")
	}

	c.Set("Content-Type", "text/html")
	component := templates.Dashboard(urls)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) ShowAddForm(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	component := templates.AddURLForm()
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) ShowEditForm(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	url, err := h.urlRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).SendString("URL not found")
	}

	c.Set("Content-Type", "text/html")
	component := templates.EditURLForm(url)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

// API Routes
func (h *Handler) CreateURL(c *fiber.Ctx) error {
	var url models.URL
	url.Name = c.FormValue("name")
	url.URL = c.FormValue("url")
	checkInterval := c.FormValue("check_interval")
	if checkInterval != "" {
		if interval, err := strconv.Atoi(checkInterval); err == nil {
			url.CheckInterval = interval
		} else {
			url.CheckInterval = 300 // default 5 minutes
		}
	} else {
		url.CheckInterval = 300 // default 5 minutes
	}

	if url.Name == "" || url.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and URL are required"})
	}

	err := h.urlRepo.Create(&url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create URL"})
	}

	// Check if it's a form submission (redirect) or API call (JSON)
	if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return c.Redirect("/")
	}
	return c.JSON(url)
}

func (h *Handler) UpdateURL(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	url, err := h.urlRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
	}

	url.Name = c.FormValue("name")
	url.URL = c.FormValue("url")
	checkInterval := c.FormValue("check_interval")
	if checkInterval != "" {
		if interval, err := strconv.Atoi(checkInterval); err == nil {
			url.CheckInterval = interval
		}
	}

	if url.Name == "" || url.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and URL are required"})
	}

	err = h.urlRepo.Update(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update URL"})
	}

	// Check if it's a form submission (redirect) or API call (JSON)
	if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return c.Redirect("/")
	}
	return c.JSON(url)
}

func (h *Handler) DeleteURL(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.urlRepo.Delete(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete URL"})
	}

	return c.JSON(fiber.Map{"message": "URL deleted successfully"})
}

func (h *Handler) CheckURL(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.checker.CheckSingleURL(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check URL"})
	}

	return c.JSON(fiber.Map{"message": "URL checked successfully"})
}

func (h *Handler) CheckAllURLs(c *fiber.Ctx) error {
	err := h.checker.CheckAllURLs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check URLs"})
	}

	return c.JSON(fiber.Map{"message": "All URLs checked successfully"})
}

func (h *Handler) GetURLs(c *fiber.Ctx) error {
	urls, err := h.urlRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch URLs"})
	}

	return c.JSON(urls)
}

func (h *Handler) GetURL(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	url, err := h.urlRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
	}

	return c.JSON(url)
}
