package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JSONResponseMiddleware memastikan semua response JSON konsisten,
// kecuali untuk file statis seperti Swagger UI
func JSONResponseMiddleware(c *fiber.Ctx) error {
	// biarkan lewat kalau ini route swagger atau asset statis
	if strings.HasPrefix(c.Path(), "/swagger") {
		return c.Next()
	}

	// jalankan handler route dulu
	err := c.Next()
	if err != nil {
		code := fiber.StatusInternalServerError
		msg := "internal server error"

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			msg = e.Message
		}

		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"message": msg,
			"data":    nil,
		})
	}

	// biarin aja kalau udah JSON
	if string(c.Response().Header.ContentType()) == fiber.MIMEApplicationJSON {
		return nil
	}

	// fallback buat non-JSON (misal text biasa)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    string(c.Response().Body()),
	})
}
