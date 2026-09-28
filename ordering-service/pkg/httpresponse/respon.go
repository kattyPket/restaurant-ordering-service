package httpresponse

import "github.com/gofiber/fiber/v3"

type Envelope struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(Envelope{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Envelope{
		Code:    code,
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}
