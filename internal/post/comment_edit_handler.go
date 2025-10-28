// internal/post/comment_edit_handler.go
package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

func RegisterCommentEditRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	r.Put("/:post_id/comments/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		commentID := c.Params("id")

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

		var comment Comment
		if err := db.First(&comment, commentID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "comment not found")
		}
		if comment.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "not your comment")
		}

		comment.Content = body.Content
		if err := db.Save(&comment).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update comment")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    comment,
		})
	})
}
