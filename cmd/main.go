package main

import (
	"email-sms-service/api"
	"email-sms-service/config"
	"email-sms-service/internal/storage"
	"email-sms-service/pkg/logger"
	"email-sms-service/pkg/queue"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	logger.InitLogger()

	// Load environment variables
	config.LoadConfig()

	storage.InitDB()

	storage.InitRedis()

	q := queue.NewQueue(storage.GetRedis())

	go q.ProcessEmailTasks()

	r := chi.NewRouter()

	api.SetupRoutes(r, q)

	logger.Log.Info("Starting API server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Log.Fatalf("Failed to start API server: %v", err)
	}
}
