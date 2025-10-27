package post

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	PostID  uint   `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
}
