package user

import (
	"fmt"
	"strings"

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
		if err := db.Where("follower_id = ? AND following_id = ?", userID, target.ID).
			Limit(1).Find(&existing).Error; err == nil && existing.ID != 0 {

			// Sudah follow → unfollow
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unfollow")
			}
			return c.JSON(fiber.Map{"following": false})
		}

		newFollow := Follow{FollowerID: userID, FollowingID: target.ID}
		if err := db.Create(&newFollow).Error; err != nil {
			if strings.Contains(err.Error(), "unique") {
				return c.JSON(fiber.Map{"following": true})
			}
			return fiber.NewError(fiber.StatusInternalServerError, "failed to follow")
		}

		// 🔔 Kirim notifikasi follow
		if target.ID != userID {
			notif := notification.Notification{
				UserID:  target.ID,
				ActorID: userID,
				Type:    "follow",
				Message: fmt.Sprintf("Kamu mendapatkan pengikut baru 👥"),
			}
			db.Create(&notif)
		}

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
			SELECT DISTINCT u.username
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
			SELECT DISTINCT u.username
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
