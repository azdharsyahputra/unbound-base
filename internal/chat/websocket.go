package chat

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"gorm.io/gorm"
	"unbound/internal/notification"
)

type BroadcastPayload struct {
	Type      string    `json:"type"` // message | status_update
	ChatID    uint      `json:"chat_id"`
	MessageID uint      `json:"message_id,omitempty"`
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content,omitempty"`
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type WebSocketHub struct {
	DB        *gorm.DB
	Rooms     map[uint]map[*websocket.Conn]uint
	Mutex     sync.Mutex
	Broadcast chan BroadcastPayload
}

func NewWebSocketHub(db *gorm.DB) *WebSocketHub {
	return &WebSocketHub{
		DB:        db,
		Rooms:     make(map[uint]map[*websocket.Conn]uint),
		Broadcast: make(chan BroadcastPayload),
	}
}

func (h *WebSocketHub) Run() {
	for {
		payload := <-h.Broadcast
		h.Mutex.Lock()
		conns := h.Rooms[payload.ChatID]
		for conn := range conns {
			data, _ := json.Marshal(payload)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				conn.Close()
				delete(conns, conn)
			}
		}
		h.Mutex.Unlock()
	}
}

func (h *WebSocketHub) HandleWebSocket(c *websocket.Conn) {
	chatID := c.Locals("chat_id").(uint)
	userID := c.Locals("user_id").(uint)

	h.Mutex.Lock()
	if _, ok := h.Rooms[chatID]; !ok {
		h.Rooms[chatID] = make(map[*websocket.Conn]uint)
	}
	h.Rooms[chatID][c] = userID
	h.Mutex.Unlock()

	log.Printf("User %d connected to chat %d", userID, chatID)

	if err := h.DB.Model(&Message{}).
		Where("chat_id = ? AND sender_id != ? AND status = ?", chatID, userID, "sent").
		Update("status", "delivered").Error; err == nil {
		h.Broadcast <- BroadcastPayload{
			Type:      "status_update",
			ChatID:    chatID,
			SenderID:  userID,
			Status:    "delivered",
			Timestamp: time.Now(),
		}
	}

	defer func() {
		h.Mutex.Lock()
		delete(h.Rooms[chatID], c)
		h.Mutex.Unlock()
		c.Close()
		log.Printf("User %d disconnected from chat %d", userID, chatID)
	}()

	for {
		_, msgData, err := c.ReadMessage()
		if err != nil {
			break
		}

		var incoming struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(msgData, &incoming); err != nil || incoming.Content == "" {
			continue
		}

		msg := Message{
			ChatID:   chatID,
			SenderID: userID,
			Content:  incoming.Content,
			Status:   "sent",
			IsRead:   false,
		}
		if err := h.DB.Create(&msg).Error; err != nil {
			log.Printf("❌ gagal simpan message: %v", err)
			continue
		}

		var chat Chat
		if err := h.DB.First(&chat, chatID).Error; err == nil {
			receiverID := chat.User1ID
			if receiverID == userID {
				receiverID = chat.User2ID
			}

			notif := notification.Notification{
				UserID:    receiverID,
				ActorID:   userID,
				Type:      "message",
				Message:   "Pesan baru diterima",
				IsRead:    false,
				CreatedAt: time.Now(),
			}
			if err := h.DB.Create(&notif).Error; err != nil {
				log.Printf("❌ gagal buat notif: %v", err)
			} else {
				log.Printf("Notif dikirim ke user %d dari %d", receiverID, userID)
			}
		}

		h.Broadcast <- BroadcastPayload{
			Type:      "message",
			ChatID:    chatID,
			MessageID: msg.ID,
			SenderID:  userID,
			Content:   msg.Content,
			Status:    msg.Status,
			Timestamp: msg.CreatedAt,
		}
	}
}

func (h *WebSocketHub) BroadcastEvent(chatID uint, payload BroadcastPayload) {
	h.Mutex.Lock()
	conns := h.Rooms[chatID]
	for conn := range conns {
		data, _ := json.Marshal(payload)
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
			delete(conns, conn)
		}
	}
	h.Mutex.Unlock()
}
