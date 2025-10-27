package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// FeedItem untuk response join User + Post
type FeedItem struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// RegisterFeedRoutes gabungkan post + user info
func RegisterFeedRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/feed")

	r.Get("/", func(c *fiber.Ctx) error {
		var results []FeedItem

		query := `
			SELECT p.id, u.username, p.content, p.created_at
			FROM posts p
			JOIN users u ON u.id = p.user_id
			ORDER BY p.created_at DESC
			LIMIT 50
		`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load feed")
		}

		return c.JSON(results)
	})
}
