package notification

import (
	"time"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`    // penerima notif
	ActorID   uint           `json:"actor_id"`   // yang melakukan aksi
	Type      string         `json:"type"`       // like, follow, comment
	PostID    *uint          `json:"post_id"`    // optional
	Message   string         `json:"message"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
