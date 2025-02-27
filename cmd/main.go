package main

import (
	"email-sms-service/config"
	"email-sms-service/internal/email"
	"email-sms-service/internal/storage"
	"email-sms-service/pkg/logger"
)

func main() {

	logger.InitLogger()

	// Load environment variables
	config.LoadConfig()

	storage.InitDB()

	storage.InitRedis()

	err := email.SendEmail("user@example.com", "Welcome!", "Thank you for signing up.")
	if err != nil {
		logger.Log.Errorf("Failed to send email: %v", err)
	}
}
