package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetDBConfig() string { 
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// get confidential information from .env file
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

	// format connection string
	connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname,
    )

    return connStr
}