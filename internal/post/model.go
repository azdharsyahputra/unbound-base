package post

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
}
