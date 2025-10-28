package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"unbound/internal/auth"
	"unbound/internal/common/db"
	"unbound/internal/common/middleware"
	"unbound/internal/post"
	"unbound/internal/search"
	"unbound/internal/user"
	"unbound/internal/notification"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	app.Use(middleware.JSONResponseMiddleware)

	database := db.Connect()
	authSvc := auth.NewAuthService(database)

	auth.RegisterRoutes(app, database, authSvc)
	user.RegisterRoutes(app, database)
	user.RegisterProfileRoutes(app, database)
	user.RegisterFollowRoutes(app, database, authSvc)
	post.RegisterRoutes(app, database, authSvc)
	post.RegisterLikeRoutes(app, database, authSvc)
	post.RegisterCommentRoutes(app, database, authSvc)
	post.RegisterFeedRoutes(app, database, authSvc)
	search.RegisterSearchRoutes(app, database)
	post.RegisterEditRoutes(app, database, authSvc)
	post.RegisterCommentEditRoutes(app, database, authSvc)
	notification.RegisterRoutes(app, database, authSvc)


	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"message": "Welcome to Unbound API v0.8",
			},
		})
	})

	log.Fatal(app.Listen(":8080"))
}
