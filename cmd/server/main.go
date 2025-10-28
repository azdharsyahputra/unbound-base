package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"unbound/internal/auth"
	"unbound/internal/common/db"
	"unbound/internal/post"
	"unbound/internal/user"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()
	database := db.Connect()
	authSvc := auth.NewAuthService(database)

	auth.RegisterRoutes(app, database, authSvc)
	user.RegisterRoutes(app, database)
	user.RegisterProfileRoutes(app, database)
	user.RegisterFollowRoutes(app, database, authSvc)
	post.RegisterRoutes(app, database, authSvc)
	post.RegisterLikeRoutes(app, database, authSvc)
	post.RegisterCommentRoutes(app, database, authSvc)
	post.RegisterFeedRoutes(app, database)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Unbound API v0.4 (Profile)"})
	})

	log.Fatal(app.Listen(":8080"))
}
