package chat

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"unbound/internal/notification"
)

type ChatService struct {
	DB *gorm.DB
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{DB: db}
}

func (s *ChatService) GetOrCreateChat(user1ID, user2ID uint) (*Chat, error) {
	var chat Chat
	err := s.DB.
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Preload("Messages").
		First(&chat).Error

	if err == gorm.ErrRecordNotFound {
		chat = Chat{User1ID: user1ID, User2ID: user2ID}
		if err := s.DB.Create(&chat).Error; err != nil {
			return nil, err
		}
		return &chat, nil
	}

	return &chat, err
}

func (s *ChatService) GetMessages(chatID uint) ([]Message, error) {
	var messages []Message
	err := s.DB.
		Where("chat_id = ?", chatID).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}

func (s *ChatService) SendMessage(chatID, senderID uint, content string) (*Message, error) {
	msg := Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
		Status:   "sent",
		IsRead:   false,
	}

	if err := s.DB.Create(&msg).Error; err != nil {
		return nil, err
	}

	var chat Chat
	if err := s.DB.First(&chat, chatID).Error; err == nil {
		var targetUser uint
		if chat.User1ID == senderID {
			targetUser = chat.User2ID
		} else {
			targetUser = chat.User1ID
		}

		notif := notification.Notification{
			UserID:  targetUser,
			ActorID: senderID,
			Type:    "message",
			Message: fmt.Sprintf("Pesan baru dari user %d", senderID),
			IsRead:  false,
		}

		if err := s.DB.Create(&notif).Error; err != nil {
			log.Printf("⚠️ gagal bikin notif: %v", err)
		}
	}

	return &msg, nil
}

func (s *ChatService) MarkMessagesAsDelivered(chatID, userID uint) error {
	return s.DB.
		Model(&Message{}).
		Where("chat_id = ? AND sender_id != ? AND status = ?", chatID, userID, "sent").
		Update("status", "delivered").Error
}

func (s *ChatService) MarkMessagesAsRead(chatID, userID uint, t time.Time) error {
	return s.DB.
		Model(&Message{}).
		Where("chat_id = ? AND sender_id != ? AND status != ?", chatID, userID, "read").
		Updates(map[string]interface{}{
			"status":   "read",
			"read_at":  t,
			"is_read":  true,
		}).Error
}
