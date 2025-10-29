package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username" example:"ajar_dev"`
	Email    string `gorm:"uniqueIndex;not null" json:"email" example:"ajar@example.com"`
	Password string `gorm:"not null" json:"password" example:"secret123"`
}
