package notification

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/notifications")

	r.Get("/", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		return GetAllNotifications(c, db)
	})

	r.Post("/read", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		return MarkAllAsRead(c, db)
	})
}

// ================================================================
// HANDLERS
// ================================================================

// GetAllNotifications godoc
// @Summary Get all notifications of the current user
// @Description Mengambil semua notifikasi user yang sedang login
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications [get]
func GetAllNotifications(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userID").(uint)

	var notifs []Notification
	if err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch notifications")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    notifs,
	})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Menandai semua notifikasi user sebagai sudah dibaca
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications/read [post]
func MarkAllAsRead(c *fiber.Ctx, db *gorm.DB) error {
	userID := c.Locals("userID").(uint)
	if err := db.Model(&Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update notifications")
	}
	return c.JSON(fiber.Map{"success": true})
}
