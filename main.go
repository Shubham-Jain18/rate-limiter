package main

import (
	"fmt"
	"os"
	"rate-limiter/limiter"
	rlredis "rate-limiter/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	// Initialize Redis
	rlredis.InitRedis()

	// Create a new Gin router
	router := gin.Default()

	// Simple health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Sliding Window Counter
	router.GET("/ping/sliding-window-counter",
		limiter.SlidingWindowCounterMiddleware(rlredis.RDB, 5, 60),
		pingHandler,
	)

	// Token Bucket (context passed as first arg)
	router.GET("/ping/token-bucket",
		limiter.TokenBucketMiddleware(rlredis.Ctx, rlredis.RDB, "token_bucket", 5, 1),
		pingHandler,
	)

	// Leaky Bucket (context passed as first arg)
	router.GET("/ping/leaky-bucket",
		limiter.LeakyBucketMiddleware(rlredis.Ctx, rlredis.RDB, "leaky_bucket", 1),
		pingHandler,
	)

	// Sliding Window Log (context passed as first arg)
	router.GET("/ping/sliding-window-log",
		limiter.SlidingWindowLogMiddleware(rlredis.Ctx, rlredis.RDB, "log_bucket", 5, 60),
		pingHandler,
	)
	fmt.Println("Server running on port", port)
	// Start server
	router.Run(":" + port)
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
