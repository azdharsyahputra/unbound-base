// internal/post/edit_handler.go
package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

func RegisterEditRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	// ✏️ Edit Post
	r.Put("/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		postID := c.Params("id")
		var body struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if body.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "content cannot be empty")
		}

		userID := c.Locals("userID").(uint)

		var p Post
		if err := db.First(&p, postID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "post not found")
		}
		if p.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "not your post")
		}

		p.Content = body.Content
		if err := db.Save(&p).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update post")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    p,
		})
	})
}
