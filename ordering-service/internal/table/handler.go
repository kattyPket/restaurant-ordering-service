package table

import (
	"crypto/subtle"
	"errors"
	"ordering-service/configs"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ordering-service/pkg/httpresponse" // adjust import path to match your pkgs/http package
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) StartSession(c fiber.Ctx) error {
	var req StartSessionRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	sess, dt, err := h.svc.StartSession(req.QRToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrTableNotFound):
			return httpresponse.Error(c, fiber.StatusNotFound, "table not found")
		case errors.Is(err, ErrTableNotUsable):
			return httpresponse.Error(c, fiber.StatusConflict, "table is unavailable")
		case errors.Is(err, ErrTableAwaitingCheckout):
			return httpresponse.Error(c, fiber.StatusConflict, "table is awaiting checkout")
		default:
			return httpresponse.Error(c, fiber.StatusInternalServerError, "failed to start session")
		}
	}

	resp := StartSessionResponse{
		SessionToken: sess.Token,
		TableID:      dt.ID.String(),
		TableNo:      dt.TableNo,
		ExpiresAt:    sess.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return httpresponse.Success(c, fiber.StatusOK, "session started", resp)
}

func (h *Handler) HandleCheckoutCompleted(c fiber.Ctx) error {
	secret := configs.Cfg.Internal.ServiceSecret
	expected := "Bearer " + secret
	if secret == "" || subtle.ConstantTimeCompare([]byte(c.Get("Authorization")), []byte(expected)) != 1 {
		return httpresponse.Error(c, fiber.StatusUnauthorized, "unauthorized internal event")
	}

	var event CheckoutCompletedEvent
	if err := c.Bind().Body(&event); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid event payload")
	}
	if event.Event != "table.checkout.completed" || event.EventID == "" {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid event")
	}
	sessionID, err := uuid.Parse(event.SessionID)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid session_id")
	}
	tableID, err := uuid.Parse(event.TableID)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid table_id")
	}
	if err := h.svc.CloseSessionForCheckout(sessionID, tableID); err != nil {
		switch {
		case errors.Is(err, ErrTableNotFound), errors.Is(err, gorm.ErrRecordNotFound):
			return httpresponse.Error(c, fiber.StatusNotFound, "table or session not found")
		case errors.Is(err, ErrSessionTableMismatch):
			return httpresponse.Error(c, fiber.StatusConflict, "session does not belong to table")
		default:
			return httpresponse.Error(c, fiber.StatusInternalServerError, "failed to close session")
		}
	}
	return httpresponse.Success(c, fiber.StatusOK, "checkout processed", nil)
}
