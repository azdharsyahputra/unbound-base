package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
	"unbound/internal/notification"
)

// RegisterFollowRoutes handles follow system
func RegisterFollowRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/users")

	// POST /users/:username/follow → follow/unfollow user
	r.Post("/:username/follow", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		targetUsername := c.Params("username")
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		var target auth.User
		if err := db.Where("username = ?", targetUsername).First(&target).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "target user not found")
		}

		if userID == target.ID {
			return fiber.NewError(fiber.StatusBadRequest, "you can't follow yourself")
		}

		var existing Follow
		err := db.Where("follower_id = ? AND following_id = ?", userID, target.ID).First(&existing).Error

		// sudah follow → unfollow
		if err == nil {
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unfollow")
			}
			return c.JSON(fiber.Map{"following": false})
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, "query error")
		}

		// belum follow → follow
		newFollow := Follow{FollowerID: userID, FollowingID: target.ID}
		if err := db.Create(&newFollow).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to follow")
		}

		// ===== 🔔 Buat notifikasi ke user yang di-follow =====
		if target.ID != userID {
			notif := notification.Notification{
				UserID:  target.ID,     // penerima notif
				ActorID: userID,        // pelaku follow
				Type:    "follow",
				Message: fmt.Sprintf("Kamu mendapatkan pengikut baru 👥"),
			}
			db.Create(&notif)
		}
		// =====================================================

		return c.JSON(fiber.Map{"following": true})
	})

	// GET /users/:username/followers
	r.Get("/:username/followers", func(c *fiber.Ctx) error {
		username := c.Params("username")

		var target auth.User
		if err := db.Where("username = ?", username).First(&target).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var followers []struct {
			Username string `json:"username"`
		}

		query := `
			SELECT u.username
			FROM follows f
			JOIN users u ON u.id = f.follower_id
			WHERE f.following_id = ?
		`
		if err := db.Raw(query, target.ID).Scan(&followers).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch followers")
		}

		return c.JSON(followers)
	})

	// GET /users/:username/following
	r.Get("/:username/following", func(c *fiber.Ctx) error {
		username := c.Params("username")

		var target auth.User
		if err := db.Where("username = ?", username).First(&target).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var following []struct {
			Username string `json:"username"`
		}

		query := `
			SELECT u.username
			FROM follows f
			JOIN users u ON u.id = f.following_id
			WHERE f.follower_id = ?
		`
		if err := db.Raw(query, target.ID).Scan(&following).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch following")
		}

		return c.JSON(following)
	})
}
