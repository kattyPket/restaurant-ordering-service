// handler.go
package order

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PlaceOrder(c fiber.Ctx) error {
	// ດຶງ TableID ຈາກ Middleware
	tableIDStr, ok := c.Locals("table_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session ເປົ່າຫວ່າງ ຫຼື ບໍ່ຖືກຕ້ອງ"})
	}
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid table ID"})
	}
	sessionIDStr, ok := c.Locals("session_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session"})
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session ID"})
	}

	var req PlaceOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ຮູບແບບຂໍ້ມູນບໍ່ຖືກຕ້ອງ"})
	}

	resp, err := h.service.PlaceOrder(tableID, sessionID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ສັ່ງອາຫານສຳເລັດ",
		"data":    resp,
	})
}

func (h *Handler) GetOrderTracking(c fiber.Ctx) error {
	// 1. ດຶງ TableID ຈາກ Middleware
	tableIDStr, ok := c.Locals("table_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session ເປົ່າຫວ່າງ ຫຼື ບໍ່ຖືກຕ້ອງ"})
	}

	// 2. ກວດສອບຮູບແບບ TableID
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "ຮູບແບບ Table ID ບໍ່ຖືກຕ້ອງ"})
	}
	sessionIDStr, ok := c.Locals("session_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session"})
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session ID"})
	}

	// 3. ດຶງ OrderID ຈາກ Params
	orderIDParam := c.Params("id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ຮູບແບບ Order ID ບໍ່ຖືກຕ້ອງ"})
	}

	// 4. ສົ່ງທັງ tableID ແລະ orderID ໄປໃຫ້ Service ຈັດການ
	order, err := h.service.GetOrderTracking(tableID, sessionID, orderID)

	if err != nil {
		// ຖ້າເປັນ ErrOrderNotFound ໃຫ້ຕອບ 404 (ປ້ອງກັນການເດົາ ID ຂອງໂຕະອື່ນ)
		if errors.Is(err, ErrOrderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ບໍ່ພົບຂໍ້ມູນອໍເດີ້"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ດຶງຂໍ້ມູນອໍເດີ້ສຳເລັດ",
		"data":    order,
	})
}
