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
		return UpdateCommentHandler(c, db)
	})
}

// ================================================================
// HANDLER
// ================================================================

// UpdateCommentHandler godoc
// @Summary Update comment content
// @Description Mengubah isi komentar milik user yang sedang login
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param id path int true "Comment ID"
// @Param data body UpdateCommentRequest true "Updated comment content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /posts/{post_id}/comments/{id} [put]
func UpdateCommentHandler(c *fiber.Ctx, db *gorm.DB) error {
	commentID := c.Params("id")

	var body UpdateCommentRequest
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
}

// ================================================================
// REQUEST SCHEMA
// ================================================================

// UpdateCommentRequest digunakan untuk endpoint PUT /posts/:post_id/comments/:id
type UpdateCommentRequest struct {
	Content string `json:"content" example:"Updated comment content"`
}
