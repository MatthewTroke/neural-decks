package bootstrap

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisOnce sync.Once

func NewRedisInstance(env *Env) *redis.Client {
	var redisClient *redis.Client

	// Get Redis configuration from environment
	redisHost := env.RedisHost

	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := env.RedisPort

	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := env.RedisPassword
	redisDB := env.RedisDB

	if redisDB == 0 {
		redisDB = 0
	}

	// Debug: Print Redis configuration
	log.Printf("🔍 Redis Config - Host: %s, Port: %s, DB: %d", redisHost, redisPort, redisDB)

	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisPassword,
			DB:       redisDB,
			PoolSize: 10,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Printf("Failed to connect to Redis: %v", err)
			return
		}

		log.Println("✅ Redis connected successfully")
		redisClient = client
	})

	return redisClient
}

func CloseRedis(client *redis.Client) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}
