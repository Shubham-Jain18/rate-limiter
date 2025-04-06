package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Exported context for Redis operations
var Ctx = context.Background()

// Exported Redis client
var RDB *redis.Client

// Initializes the Redis client and assigns it to RDB
func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adjust if needed
		Password: "",               // Add password if set
		DB:       0,                // Default DB
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis successfully")
}
