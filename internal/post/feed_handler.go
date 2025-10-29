package post

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

type FeedItem struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Likes     int64  `json:"likes"`
}

func RegisterFeedRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	r := app.Group("/feed")

	r.Get("/", func(c *fiber.Ctx) error {
		var results []FeedItem

		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		sortOrder := c.Query("sort", "newest")

		if limit <= 0 || limit > 100 {
			limit = 20
		}

		order := "DESC"
		if sortOrder == "oldest" {
			order = "ASC"
		}

		query := `
			SELECT p.id, u.username, p.content, p.created_at,
				COUNT(DISTINCT l.id) AS likes
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN likes l ON l.post_id = p.id
			GROUP BY p.id, u.username, p.content, p.created_at
			ORDER BY p.created_at ` + order + `
			LIMIT ? OFFSET ?
		`

		if err := db.Raw(query, limit, offset).Scan(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load feed")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    results,
			"meta": fiber.Map{
				"limit":  limit,
				"offset": offset,
				"sort":   sortOrder,
				"count":  len(results),
			},
		})
	})

	r.Get("/following", middleware.JWTProtected(authSvc), func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		var results []FeedItem

		limit, _ := strconv.Atoi(c.Query("limit", "20"))
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		sortOrder := c.Query("sort", "newest")

		if limit <= 0 || limit > 100 {
			limit = 20
		}

		order := "DESC"
		if sortOrder == "oldest" {
			order = "ASC"
		}

		query := `
			SELECT p.id, u.username, p.content, p.created_at,
				COUNT(DISTINCT l.id) AS likes
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN likes l ON l.post_id = p.id
			WHERE p.user_id IN (
				SELECT following_id FROM follows WHERE follower_id = ?
			)
			OR p.user_id = ?
			GROUP BY p.id, u.username, p.content, p.created_at
			ORDER BY p.created_at ` + order + `
			LIMIT ? OFFSET ?
		`

		if err := db.Raw(query, userID, userID, limit, offset).Scan(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to load following feed")
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    results,
			"meta": fiber.Map{
				"limit":  limit,
				"offset": offset,
				"sort":   sortOrder,
				"count":  len(results),
			},
		})
	})
}
