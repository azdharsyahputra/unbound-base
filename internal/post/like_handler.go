package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
	"unbound/internal/common/utils"
)


// RegisterLikeRoutes handles /posts/:id/like and /posts/:id/likes
func RegisterLikeRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	// POST /posts/:id/like → toggle like/unlike
	r.Post("/:id/like", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		postID := c.Params("id")
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		var existing Like
		err := db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existing).Error
		if err == nil {
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unlike")
			}
			return c.JSON(fiber.Map{"liked": false})
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, "query error")
		}

		newLike := Like{UserID: userID, PostID: utils.ToUint(postID)}
		if err := db.Create(&newLike).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to like")
		}

		return c.JSON(fiber.Map{"liked": true})
	})

	// GET /posts/:id/likes → total like
	r.Get("/:id/likes", func(c *fiber.Ctx) error {
		postID := c.Params("id")
		var count int64
		if err := db.Model(&Like{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to count likes")
		}
		return c.JSON(fiber.Map{"post_id": postID, "likes": count})
	})
}
