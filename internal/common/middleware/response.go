package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func JSONResponseMiddleware(c *fiber.Ctx) error {u
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

	if string(c.Response().Header.ContentType()) == fiber.MIMEApplicationJSON {
		return nil
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    string(c.Response().Body()),
	})
}
