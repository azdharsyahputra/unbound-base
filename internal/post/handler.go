package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

type createPostReq struct {
	Content string `json:"content"`
}

func RegisterRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	r.Get("/", func(c *fiber.Ctx) error {
		var posts []Post
		if err := db.Order("id DESC").Limit(100).Find(&posts).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch posts")
		}
		return c.JSON(posts)
	})

	r.Post("/", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		var req createPostReq
		if err := c.BodyParser(&req); err != nil || req.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "content is required")
		}
		userID, ok := c.Locals("userID").(uint)
		if !ok || userID == 0 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}
		p := &Post{
			UserID:  userID,
			Content: req.Content,
		}
		if err := db.Create(p).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create post")
		}
		return c.Status(fiber.StatusCreated).JSON(p)
	})
}
