package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// JSONResponseMiddleware memastikan semua response JSON konsisten
func JSONResponseMiddleware(c *fiber.Ctx) error {
	// jalankan handler route dulu
	err := c.Next()

	// jika ada error dari handler → ubah ke format JSON standar
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

	// pastikan kalau response bukan JSON, ubah jadi format standar
	if string(c.Response().Header.ContentType()) == fiber.MIMEApplicationJSON {
		return nil // udah JSON, biarin aja
	}

	// fallback untuk non-JSON response (misal string atau plain text)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    string(c.Response().Body()),
	})
}
