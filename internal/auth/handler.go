package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/auth")

	r.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Register endpoint"})
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Login endpoint"})
	})
}
