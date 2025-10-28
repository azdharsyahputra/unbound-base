package search

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SearchResult struct {
	Type      string `json:"type"`
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func RegisterSearchRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/search")

	r.Get("/", func(c *fiber.Ctx) error {
		query := c.Query("query")
		filterType := c.Query("type")
		sortOrder := c.Query("sort")

		if query == "" {
			return fiber.NewError(fiber.StatusBadRequest, "query parameter is required")
		}

		var results []SearchResult
		pattern := "%" + query + "%"

		switch filterType {
		case "user":
			sql := `
				SELECT 'user' AS type, id, username AS content, NULL AS created_at
				FROM users
				WHERE username ILIKE ?
				LIMIT 50
			`
			if err := db.Raw(sql, pattern).Scan(&results).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to search users")
			}

		case "post":
			order := "DESC"
			if sortOrder == "oldest" {
				order = "ASC"
			}

			sql := `
				SELECT 'post' AS type, id, content, created_at
				FROM posts
				WHERE content ILIKE ?
				ORDER BY created_at ` + order + `
				LIMIT 50
			`
			if err := db.Raw(sql, pattern).Scan(&results).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to search posts")
			}

		default:
			sql := `
				SELECT 'user' AS type, id, username AS content, NULL AS created_at FROM users WHERE username ILIKE ?
				UNION
				SELECT 'post' AS type, id, content, created_at FROM posts WHERE content ILIKE ?
				LIMIT 50
			`
			if err := db.Raw(sql, pattern, pattern).Scan(&results).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to perform search")
			}
		}

		return c.JSON(results)
	})
}
