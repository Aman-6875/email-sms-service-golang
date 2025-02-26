package main

import (
	"email-sms-service/config"
	"email-sms-service/internal/storage"
	"email-sms-service/pkg/logger"
	"email-sms-service/pkg/queue"
	"fmt"
)

func main() {

	logger.InitLogger()

	// Load environment variables
	config.LoadConfig()

	storage.InitDB()

	storage.InitRedis()

	q := queue.NewQueue(storage.GetRedis())

	message := map[string]string{
		"type":    "email",
		"to":      "user@example.com",
		"subject": "Welcome!",
		"body":    "Thank you for signing up.",
	}

	err := q.Enqueue("email_queue", message)
	if err != nil {
		logger.Log.Errorf("Failed to enqueue message: %v", err)
	}

	go func() {
		err := q.Dequeue("email_queue", func(message string) {
			logger.Log.Infof("Processing message: %s", message)
		})

		if err != nil {
			logger.Log.Errorf("Failed to dequeue message: %v", err)
		}
	}()

	fmt.Println("Application is running. Press Ctrl+C to exit.")
	select {}
}
