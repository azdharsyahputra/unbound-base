package notification

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/notifications")

	// 🔔 Ambil semua notifikasi user
	r.Get("/", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)

		var notifs []Notification
		if err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch notifications")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    notifs,
		})
	})

	// ✅ Tandai semua sebagai read
	r.Post("/read", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		if err := db.Model(&Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update notifications")
		}
		return c.JSON(fiber.Map{"success": true})
	})
}
