package storage

import (
	"context"
	"email-sms-service/config"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisHost := config.GetEnv("REDIS_HOST", "localhost")
	redisPort := config.GetEnv("REDIS_PORT", "6379")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	log.Println("Successfully Connected to Redis!")
}

func GetRedis() *redis.Client {
	return RedisClient
}
