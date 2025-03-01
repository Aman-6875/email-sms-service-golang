package main

import (
	"email-sms-service/config"
	"email-sms-service/internal/email"
	"email-sms-service/internal/storage"
	"email-sms-service/pkg/logger"
	"email-sms-service/pkg/queue"
)

func main() {

	logger.InitLogger()

	// Load environment variables
	config.LoadConfig()

	storage.InitDB()

	storage.InitRedis()

	q := queue.NewQueue(storage.GetRedis())

	task := email.EmailTask{
		To:      "user@example.com",
		Subject: "Welcome!",
		Body:    "Thank you for signing up.",
	}

	if err := q.EnqueueEmailTask(task); err != nil {
		logger.Log.Errorf("Failed to enqueue email task: %v", err)
	}

	go q.ProcessEmailTasks()

	logger.Log.Info("Application is running. Press Ctrl+C to exit.")
	select {}
}
