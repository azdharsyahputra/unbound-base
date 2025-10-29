package post

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
	"unbound/internal/common/utils"
	"unbound/internal/notification"
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
		if err := db.Where("user_id = ? AND post_id = ?", userID, postID).
			Limit(1).Find(&existing).Error; err == nil && existing.ID != 0 {

			// sudah like → unlike
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unlike")
			}
			return c.JSON(fiber.Map{"liked": false})
		}

		newLike := Like{UserID: userID, PostID: utils.ToUint(postID)}
		if err := db.Create(&newLike).Error; err != nil {
			if strings.Contains(err.Error(), "unique") {
				return c.JSON(fiber.Map{"liked": true})
			}
			return fiber.NewError(fiber.StatusInternalServerError, "failed to like")
		}

		// 🔍 Ambil data post + pemiliknya
		var postOwner struct {
			ID        uint
			OwnerName string
			ActorName string
		}
		db.Raw(`
			SELECT p.user_id AS id, 
			       u.username AS owner_name,
			       (SELECT username FROM users WHERE id = ?) AS actor_name
			FROM posts p
			JOIN users u ON u.id = p.user_id
			WHERE p.id = ?
		`, userID, postID).Scan(&postOwner)

		// 🔔 Notifikasi ke pemilik post (jika bukan dirinya sendiri)
		if postOwner.ID != userID && postOwner.ID != 0 {
			notif := notification.Notification{
				UserID:  postOwner.ID,       // penerima notif
				ActorID: userID,             // pelaku like
				Type:    "like",
				PostID:  utils.ToUintPtr(postID),
				Message: fmt.Sprintf("%s menyukai postinganmu ❤️", postOwner.ActorName),
			}
			db.Create(&notif)
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
