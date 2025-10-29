package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

// RegisterEditRoutes handles /posts/:id (PUT)
func RegisterEditRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	r.Put("/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		return UpdatePostHandler(c, db)
	})
}

// ================================================================
// HANDLER
// ================================================================

// UpdatePostHandler godoc
// @Summary Update an existing post
// @Description Mengubah konten postingan milik user yang sedang login
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param data body UpdatePostRequest true "Updated post content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [put]
func UpdatePostHandler(c *fiber.Ctx, db *gorm.DB) error {
	postID := c.Params("id")
	var body UpdatePostRequest
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
}

// ================================================================
// REQUEST SCHEMA
// ================================================================

// UpdatePostRequest digunakan untuk endpoint PUT /posts/:id
type UpdatePostRequest struct {
	Content string `json:"content" example:"Updated caption from Unbound!"`
}
