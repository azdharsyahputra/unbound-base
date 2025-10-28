package post

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

// FeedItem untuk response join User + Post
type FeedItem struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Likes     int64  `json:"likes"`
}

// RegisterFeedRoutes gabungkan post + user info
func RegisterFeedRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/feed")

	// 🌍 Public Feed (semua post)
	r.Get("/", func(c *fiber.Ctx) error {
		var results []FeedItem

		query := `
			SELECT p.id, u.username, p.content, p.created_at,
				COUNT(l.id) AS likes
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN likes l ON l.post_id = p.id
			GROUP BY p.id, u.username, p.content, p.created_at
			ORDER BY p.created_at DESC
			LIMIT 50
		`

		if err := db.Raw(query).Scan(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load feed")
		}

		return c.JSON(results)
	})

	// 👥 Following Feed (hanya user yg diikuti)
	r.Get("/following", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		var results []FeedItem

		query := `
			SELECT p.id, u.username, p.content, p.created_at,
				COUNT(l.id) AS likes
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN likes l ON l.post_id = p.id
			WHERE p.user_id IN (
				SELECT following_id FROM follows WHERE follower_id = ?
			)
			OR p.user_id = ?
			GROUP BY p.id, u.username, p.content, p.created_at
			ORDER BY p.created_at DESC
			LIMIT 50
		`

		if err := db.Raw(query, userID, userID).Scan(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load following feed")
		}

		return c.JSON(results)
	})
}
