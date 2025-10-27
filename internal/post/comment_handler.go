package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
	"unbound/internal/common/utils"
)


// RegisterCommentRoutes handles /posts/:id/comments
func RegisterCommentRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	// POST /posts/:id/comments
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

		return c.Status(fiber.StatusCreated).JSON(comment)
	})

	// GET /posts/:id/comments
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
}
