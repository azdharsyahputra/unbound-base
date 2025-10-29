package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

// Request body untuk create/edit post
type createPostReq struct {
	Content string `json:"content"`
}

// RegisterRoutes godoc
// @Summary Register post routes
// @Tags Posts
func RegisterRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	// ============================================================
	// @Summary Get all posts
	// @Description Mengambil semua postingan (max 100 terbaru)
	// @Tags Posts
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts [get]
	// ============================================================
	r.Get("/", func(c *fiber.Ctx) error {
		var posts []Post
		if err := db.Order("id DESC").Limit(100).Find(&posts).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch posts")
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    posts,
		})
	})

	// ============================================================
	// @Summary Create a new post
	// @Description Membuat postingan baru untuk user yang sedang login
	// @Tags Posts
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param data body createPostReq true "Post content"
	// @Success 201 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts [post]
	// ============================================================
	r.Post("/", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		var req createPostReq
		if err := c.BodyParser(&req); err != nil || req.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "content is required")
		}
		userID, ok := c.Locals("userID").(uint)
		if !ok || userID == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		p := &Post{UserID: userID, Content: req.Content}
		if err := db.Create(p).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create post")
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    p,
		})
	})

	// ============================================================
	// @Summary Update an existing post
	// @Description Edit postingan milik sendiri
	// @Tags Posts
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param id path int true "Post ID"
	// @Param data body createPostReq true "Updated content"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts/{id} [put]
	// ============================================================
	r.Put("/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req createPostReq
		if err := c.BodyParser(&req); err != nil || req.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}

		userID := c.Locals("userID").(uint)
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "post not found")
		}
		if post.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "not your post")
		}

		post.Content = req.Content
		if err := db.Save(&post).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update post")
		}

		return c.JSON(fiber.Map{"success": true, "data": post})
	})

	// ============================================================
	// @Summary Delete a post
	// @Description Hapus postingan milik sendiri
	// @Tags Posts
	// @Security BearerAuth
	// @Param id path int true "Post ID"
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Failure 404 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts/{id} [delete]
	// ============================================================
	r.Delete("/:id", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(uint)

		var post Post
		if err := db.First(&post, id).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "post not found")
		}
		if post.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "not your post")
		}

		if err := db.Delete(&post).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to delete post")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "post deleted successfully",
		})
	})
}
