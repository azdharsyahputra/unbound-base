package notification

import (
	"time"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	ActorID   uint           `json:"actor_id"`
	Type      string         `json:"type"`
	PostID    *uint          `json:"post_id"`
	Message   string         `json:"message"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
