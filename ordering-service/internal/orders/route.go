// routes.go
package order

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, handler *Handler, sessionGuard fiber.Handler) {
	orderGroup := router.Group("/order", sessionGuard)

	orderGroup.Post("/place", handler.PlaceOrder)
	orderGroup.Get("/:id", handler.GetOrderTracking) // ສໍາລັບລູກຄ້າ track ສະຖານະ
}
