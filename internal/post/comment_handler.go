package post

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
	"unbound/internal/common/utils"
	"unbound/internal/notification"
)

func RegisterCommentRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	r.Post("/:id/comments", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		postID := c.Params("id")
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		var body struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&body); err != nil || body.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "content is required")
		}

		comment := Comment{
			UserID:  userID,
			PostID:  utils.ToUint(postID),
			Content: body.Content,
		}

		if err := db.Create(&comment).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create comment")
		}

		var postOwner struct {
			ID       uint
			Username string
		}
		if err := db.Raw(`
			SELECT p.user_id AS id, u.username 
			FROM posts p 
			JOIN users u ON u.id = p.user_id 
			WHERE p.id = ?
		`, postID).Scan(&postOwner).Error; err == nil && postOwner.ID != userID {
			notif := notification.Notification{
				UserID:  postOwner.ID,
				ActorID: userID,
				Type:    "comment",
				PostID:  utils.ToUintPtr(postID),
				Message: fmt.Sprintf("%s mengomentari postinganmu", postOwner.Username),
			}
			db.Create(&notif)
		}

		return c.Status(fiber.StatusCreated).JSON(comment)
	})

	r.Get("/:id/comments", func(c *fiber.Ctx) error {
		postID := c.Params("id")
		var comments []struct {
			ID        uint   `json:"id"`
			Username  string `json:"username"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
		}

		query := `
			SELECT c.id, u.username, c.content, c.created_at
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.post_id = ?
			ORDER BY c.created_at ASC
		`
		if err := db.Raw(query, postID).Scan(&comments).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch comments")
		}

		return c.JSON(comments)
	})

	r.Delete("/:post_id/comments/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		commentID := c.Params("id")
		userID := c.Locals("userID").(uint)

		var comment Comment
		if err := db.First(&comment, commentID).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "comment not found")
		}

		if comment.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "cannot delete others' comment")
		}

		if err := db.Delete(&comment).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to delete comment")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "comment deleted successfully",
		})
	})
}
