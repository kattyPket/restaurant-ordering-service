package table

import "github.com/gofiber/fiber/v3"

// SetupTableRoutes registers all /table endpoints on the given router group.
func SetupTableRoutes(router fiber.Router, h *Handler) {
	router.Post("/table/session", h.StartSession)
}
