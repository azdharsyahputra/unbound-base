package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"unbound/internal/auth"
	"unbound/internal/post"
	"unbound/internal/user"
)

// Connect membuka koneksi ke PostgreSQL dan menjalankan migrasi model
func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Migrasi semua tabel utama
	err = db.AutoMigrate(
		&auth.User{},
		&post.Post{},
		&post.Like{},
		&post.Comment{},
		&user.Follow{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database connected & migrated successfully")
	return db
}
