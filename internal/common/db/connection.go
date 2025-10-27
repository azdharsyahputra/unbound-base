package db

import (
	"fmt"
	"log"
	"os"
	"unbound/internal/auth"
	"unbound/internal/post"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = db.AutoMigrate(&auth.User{}, &post.Post{})
	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	log.Println("Database connected & migrated successfully")
	return db
}
