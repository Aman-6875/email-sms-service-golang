package api

import (
	"email-sms-service/pkg/queue"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, q *queue.Queue) {
	r.Post("/email", SendEmailHandler(q))
}