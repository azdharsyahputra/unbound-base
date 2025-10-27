package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, svc *AuthService) {
	r := app.Group("/auth")

	r.Post("/register", func(c *fiber.Ctx) error {
		var req RegisterReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		u, err := svc.Register(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		})
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		var req LoginReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		tok, err := svc.Login(req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return c.JSON(tok)
	})
}
