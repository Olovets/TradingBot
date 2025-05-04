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
	// Construct the DSN with password from an environment variable
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		"209.97.178.26",    // PostgreSQL server IP address
		"canadianvisa-crm", // User from environment variable
		"yzwrfr0ycacwq9ne", // Password from environment variable
		"canadianvisa-crm", // Database name from environment variable
		"5432",             // PostgreSQL default port
		sslMode,            // SSL Mode
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected successfully!")
	return db, nil
}
