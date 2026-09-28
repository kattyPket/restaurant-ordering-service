package menus

import "github.com/gofiber/fiber/v3"

func SetupRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/menu")
	api.Get("/", h.GetMenu)
}
