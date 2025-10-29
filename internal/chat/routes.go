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

// RegisterChatRoutes mendaftarkan semua route untuk fitur chat (REST + WebSocket)
func RegisterChatRoutes(app *fiber.App, db *gorm.DB, authSvc *auth.AuthService) {
	// Inisialisasi service dan hub
	svc := NewChatService(db)
	hub := NewWebSocketHub(db)
	go hub.Run() // jalankan listener broadcast

	// Handler sekarang butuh hub untuk broadcast status read
	h := NewChatHandler(svc, hub)

	// ===== REST API Chat =====
	// Gunakan JWTProtected untuk validasi token dari header Authorization
	r := app.Group("/chats", middleware.JWTProtected(authSvc))

	// Ambil daftar chat milik user
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

	// Buat atau ambil chat antar user
	r.Post("/:user_id", h.GetOrCreateChat)

	// Ambil semua pesan di chat
	r.Get("/:chat_id/messages", h.GetMessages)

	// Kirim pesan baru
	r.Post("/:chat_id/messages", h.SendMessage)

	// Tandai semua pesan sebagai 'read'
	r.Put("/:chat_id/read", h.MarkAsRead)

	// ===== WebSocket Chat =====
	// Gunakan query token: ?token=<JWT>
	app.Get("/ws/chat/:chat_id",
		middleware.WebSocketAuth(os.Getenv("JWT_SECRET")),
		websocket.New(func(c *websocket.Conn) {
			chatID, _ := strconv.Atoi(c.Params("chat_id"))
			c.Locals("chat_id", uint(chatID))
			hub.HandleWebSocket(c)
		}),
	)
}
