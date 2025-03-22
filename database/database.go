package database

import (
	"fmt"
	"log"

	"github.com/dockrelix/dockrelix-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "database/database.db"
	fmt.Println("Connecting to SQLite database...")

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to SQLite database.")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.StackDraft{},
	)
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}
}

// For tests only - creates an in-memory database
func InitDBForTesting() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	return DB
}
