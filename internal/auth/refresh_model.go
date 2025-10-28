package auth

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint
	Token     string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}
