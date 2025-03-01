package queue

import (
	"context"
	"email-sms-service/internal/email"
	"email-sms-service/pkg/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	client *redis.Client
}

func NewQueue(client *redis.Client) *Queue {
	return &Queue{
		client: client,
	}
}

func (q *Queue) Enqueue(queueName string, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	messageJSON, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("failed to serialized message: %v", err)
	}

	err = q.client.LPush(ctx, queueName, messageJSON).Err()

	if err != nil {
		return fmt.Errorf("failed to enqueue message: %v", err)
	}

	logger.Log.Infof("Message enqueued successfully in queue: %s", queueName)

	return nil
}

func (q *Queue) Dequeue(queueName string, handler func(message string)) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Blocking pop from the Redis list
	result, err := q.client.BRPop(ctx, 0*time.Second, queueName).Result()
	if err != nil {
		return fmt.Errorf("failed to dequeue message: %v", err)
	}

	// Process the message
	message := result[1]
	handler(message)

	logger.Log.Infof("Message dequeued and processed from queue: %s", queueName)
	return nil
}

func (q *Queue) EnqueueEmailTask(task email.EmailTask) error {
	return q.Enqueue("email_queue", task)
}

func (q *Queue) ProcessEmailTasks() {
	for {
		err := q.Dequeue("email_queue", func(message string) {
			var task email.EmailTask
			if err := json.Unmarshal([]byte(message), &task); err != nil {
				logger.Log.Errorf("failed to deserialize email task: %v", err)
				return
			}

			if err := email.SendEmail(task.To, task.Subject, task.Body); err != nil {
				logger.Log.Errorf("Failed to send email after retries: %v", err)

				if err := q.Enqueue("email_queue_dlq", task); err != nil {
                    logger.Log.Errorf("Failed to move task to DLQ: %v", err)
                }
				return
			}

		})

		if err != nil {
			logger.Log.Errorf("failed to dequeue email task: %v", err)
		}
	}
}
