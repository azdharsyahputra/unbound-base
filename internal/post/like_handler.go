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

// RegisterLikeRoutes godoc
// @Summary Like & Unlike system
// @Tags Likes
func RegisterLikeRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/posts")

	// LikePostHandler godoc
	// @Summary Like or Unlike a post
	// @Description Toggle like pada postingan (jika sudah like maka unlike)
	// @Tags Likes
	// @Security BearerAuth
	// @Param id path int true "Post ID"
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts/{id}/like [post]
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
				Type:    "like",
				PostID:  utils.ToUintPtr(postID),
				Message: fmt.Sprintf("%s menyukai postinganmu ❤️", postOwner.Username),
			}
			db.Create(&notif)
		}

		return c.JSON(fiber.Map{"liked": true})
	})

	// GetLikesHandler godoc
	// @Summary Get like count
	// @Description Mengambil jumlah like pada postingan
	// @Tags Likes
	// @Param id path int true "Post ID"
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /posts/{id}/likes [get]
	r.Get("/:id/likes", func(c *fiber.Ctx) error {
		postID := c.Params("id")
		var count int64
		if err := db.Model(&Like{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to count likes")
		}
		return c.JSON(fiber.Map{"post_id": postID, "likes": count})
	})
}
