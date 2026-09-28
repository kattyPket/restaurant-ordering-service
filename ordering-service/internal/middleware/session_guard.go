package middleware

import (
	"ordering-service/internal/table"
	"ordering-service/pkg/httpresponse"

	"github.com/gofiber/fiber/v3"
)

func TableSessionGuard(svc *table.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Session-Token")
		if token == "" {
			return httpresponse.Error(c, fiber.StatusUnauthorized, "missing session token")
		}

		sess, err := svc.ValidationToken(token)
		if err != nil {
			return httpresponse.Error(c, fiber.StatusUnauthorized, "invalid or expired session")
		}

		c.Locals("table_id", sess.TableID.String())
		c.Locals("session_id", sess.ID.String())
		return c.Next()
	}
}
