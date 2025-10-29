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

// RegisterCommentRoutes handles /posts/:id/comments
func RegisterCommentRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	r.Post("/:id/comments", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		return CreateCommentHandler(c, db)
	})

	r.Get("/:id/comments", func(c *fiber.Ctx) error {
		return GetCommentsHandler(c, db)
	})
}

// ================================================================
// HANDLERS
// ================================================================

// CreateCommentHandler godoc
// @Summary Create a comment on a post
// @Description Membuat komentar baru pada postingan tertentu (hanya untuk user login)
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param data body CreateCommentRequest true "Comment content"
// @Success 201 {object} Comment
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/comments [post]
func CreateCommentHandler(c *fiber.Ctx, db *gorm.DB) error {
	postID := c.Params("id")
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	var body CreateCommentRequest
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

	// 🔔 Kirim notifikasi ke pemilik post
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
			Message: fmt.Sprintf("%s mengomentari postinganmu 💬", postOwner.Username),
		}
		db.Create(&notif)
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

// GetCommentsHandler godoc
// @Summary Get all comments for a post
// @Description Mengambil semua komentar dari sebuah postingan
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} CommentResponse
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/comments [get]
func GetCommentsHandler(c *fiber.Ctx, db *gorm.DB) error {
	postID := c.Params("id")
	var comments []CommentResponse

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
}

// ================================================================
// REQUEST & RESPONSE SCHEMA
// ================================================================

// CreateCommentRequest digunakan untuk endpoint POST /posts/:id/comments
type CreateCommentRequest struct {
	Content string `json:"content" example:"Komentar yang menarik banget!"`
}

// CommentResponse digunakan untuk menampilkan komentar dari GET /posts/:id/comments
type CommentResponse struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"ajar_dev"`
	Content   string `json:"content" example:"Mantap banget postingannya!"`
	CreatedAt string `json:"created_at" example:"2025-10-29T09:00:00Z"`
}
