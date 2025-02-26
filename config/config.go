package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// LoadConfig loads environment variables from the .env file
func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

// GetEnv retrieves the value of an environment variable
func GetEnv(key string, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}