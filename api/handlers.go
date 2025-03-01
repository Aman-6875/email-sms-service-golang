package api

import (
	"email-sms-service/internal/email"
	"email-sms-service/pkg/queue"
	"encoding/json"
	"net/http"
	"strings"
)

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmailHandler (q *queue.Queue) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var req SendEmailRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}


		if !isValidEmail(req.To) {
			respondWithError(w, http.StatusBadRequest, "Invalid email address")
			return
		}

		task := email.EmailTask{
			To:      req.To,
			Subject: req.Subject,
			Body:    req.Body,
		}

		if err := q.EnqueueEmailTask(task); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to enqueue email task")
			return
		}

		respondWithJSON(w, http.StatusAccepted, map[string]string{"message": "Email task enqueued successfully"})

	}
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}


func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func isValidEmail(email string) bool {
	// Basic email validation (you can use a library for more robust validation)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}