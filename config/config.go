package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	// Load environment variables from ..env
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading ..env file")
	}
	sslMode := os.Getenv("DB_SSLMODE")

	// Construct the database connection URL
	dsn := fmt.Sprintf("host=%s user=%s  dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		sslMode,
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected successfully!")
	return db, nil
}
