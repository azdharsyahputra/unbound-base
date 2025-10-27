package main

import (
	"log"
	"unbound/internal/auth"
	"unbound/internal/common/db"
	"unbound/internal/post"
	"unbound/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	_ = godotenv.Load()

	app := fiber.New()

	// Connect DB
	database := db.ConnectDB()

	// Register routes
	auth.RegisterRoutes(app, database)
	user.RegisterRoutes(app, database)
	post.RegisterRoutes(app, database)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Unbound API v0.1"})
	})

	log.Fatal(app.Listen(":8080"))
}
