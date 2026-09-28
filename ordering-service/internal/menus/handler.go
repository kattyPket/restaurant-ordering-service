package menus

import (
	"ordering-service/pkg/httpresponse"

	"github.com/gofiber/fiber/v3"
)

type Handler struct{ service Service }

func NewHandler(s Service) *Handler { return &Handler{s} }

func (h *Handler) GetMenu(c fiber.Ctx) error {
	menu, err := h.service.GetMenu()
	if err != nil {
		return httpresponse.Error(c, fiber.StatusInternalServerError, "unable to retrieve menu")
	}

	return httpresponse.Success(c, fiber.StatusOK, "menu retrieved successfully", menu)
}
