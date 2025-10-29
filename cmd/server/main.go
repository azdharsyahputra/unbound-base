//go:generate swag init -g cmd/server/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	// import docs hasil swag init
	_ "unbound/docs"

	"unbound/internal/auth"
	"unbound/internal/common/db"
	"unbound/internal/common/middleware"
	"unbound/internal/notification"
	"unbound/internal/post"
	"unbound/internal/search"
	"unbound/internal/user"
)

// @title Unbound API
// @version 1.0
// @description REST API backend for Unbound microservice social platform.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email unbound@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()

	app := fiber.New()

	// Middleware global, tapi skip /swagger biar UI gak keubah ke JSON
	app.Use(func(c *fiber.Ctx) error {
		if len(c.Path()) >= 8 && c.Path()[:8] == "/swagger" {
			return c.Next()
		}
		return middleware.JSONResponseMiddleware(c)
	})

	database := db.Connect()
	authSvc := auth.NewAuthService(database)

	// Register semua routes
	auth.RegisterRoutes(app, database, authSvc)
	user.RegisterRoutes(app, database)
	user.RegisterProfileRoutes(app, database)
	user.RegisterFollowRoutes(app, database, authSvc)
	post.RegisterRoutes(app, database, authSvc)
	post.RegisterLikeRoutes(app, database, authSvc)
	post.RegisterCommentRoutes(app, database, authSvc)
	post.RegisterFeedRoutes(app, database, authSvc)
	post.RegisterEditRoutes(app, database, authSvc)
	post.RegisterCommentEditRoutes(app, database, authSvc)
	search.RegisterSearchRoutes(app, database)
	notification.RegisterRoutes(app, database, authSvc)

	// ✅ Swagger route (otomatis handle semua file)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"message": "Welcome to Unbound API v1.0",
			},
		})
	})

	log.Fatal(app.Listen(":8080"))
}
