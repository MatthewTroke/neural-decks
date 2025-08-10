package infra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInstanceArgs struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func NewRedisInstanceArgs(host string, port string, password string, db int) RedisInstanceArgs {
	return RedisInstanceArgs{
		RedisHost:     host,
		RedisPort:     port,
		RedisPassword: password,
		RedisDB:       db,
	}
}

func NewRedisInstance(args RedisInstanceArgs) *redis.Client {
	var redisClient *redis.Client

	redisHost := args.RedisHost
	redisPort := args.RedisPort
	redisPassword := args.RedisPassword
	redisDB := args.RedisDB

	log.Printf("🔍 Redis Config - Host: %s, Port: %s, DB: %d", redisHost, redisPort, redisDB)

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
		return nil
	}

	log.Println("✅ Redis connected successfully")

	redisClient = client

	return redisClient
}

func CloseRedis(client *redis.Client) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}
