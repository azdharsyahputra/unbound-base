package chat

import (
	"time"
)

type Chat struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	User1ID   uint       `json:"user1_id"`
	User2ID   uint       `json:"user2_id"`
	CreatedAt time.Time  `json:"created_at"`
	Messages  []Message  `json:"messages" gorm:"foreignKey:ChatID"`
}

type Message struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ChatID     uint       `json:"chat_id"`
	SenderID   uint       `json:"sender_id"`
	Content    string     `json:"content"`

	Status     string     `gorm:"type:varchar(20);default:'sent'" json:"status"`
	IsRead     bool       `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time `json:"read_at,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
}
