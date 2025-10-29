package chat

import (
	"os"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unbound/internal/auth"
	"unbound/internal/common/middleware"
)

func RegisterChatRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	// Inisialisasi service dan hub
	svc := NewChatService(db)
	hub := NewWebSocketHub(db)
	go hub.Run()

	h := NewChatHandler(svc, hub)

	r := app.Group("/chats", middleware.JWTProtected(authSvc))

	r.Get("/", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		var chats []Chat

		if err := db.
			Where("user1_id = ? OR user2_id = ?", userID, userID).
			Preload("Messages", func(db *gorm.DB) *gorm.DB {
				return db.Order("messages.created_at desc").Limit(1)
			}).
			Find(&chats).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(chats)
	})

	r.Post("/:user_id", h.GetOrCreateChat)
	r.Get("/:chat_id/messages", h.GetMessages)
	r.Post("/:chat_id/messages", h.SendMessage)
	r.Put("/:chat_id/read", h.MarkAsRead)

	app.Get("/ws/chat/:chat_id",
		middleware.WebSocketAuth(os.Getenv("JWT_SECRET")),
		websocket.New(func(c *websocket.Conn) {
			chatID, _ := strconv.Atoi(c.Params("chat_id"))
			c.Locals("chat_id", uint(chatID))
			hub.HandleWebSocket(c)
		}),
	)
}
