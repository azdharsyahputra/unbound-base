package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/posts")

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List all posts"})
	})

	r.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Create new post"})
	})
}
