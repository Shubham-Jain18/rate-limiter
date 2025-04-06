package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func LeakyBucketMiddleware(ctx context.Context, rdb *redis.Client, key string, leakRate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		lastLeakKey := key + ":last_leak"
		waterKey := key + ":water"

		// Set maximum bucket capacity (stricter)
		maxCapacity := 1

		// Get last leak time and current water level
		lastLeak, err := rdb.Get(ctx, lastLeakKey).Int64()
		if err != nil {
			lastLeak = now
		}
		water, err := rdb.Get(ctx, waterKey).Int()
		if err != nil {
			water = 0
		}

		// Leak water based on elapsed time
		elapsed := now - lastLeak
		leaked := int(elapsed) * leakRate
		water = max(0, water-leaked)

		if water >= maxCapacity {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded (Leaky Bucket)"})
			return
		}

		// Update water level and last leak time
		pipe := rdb.TxPipeline()
		pipe.Set(ctx, waterKey, water+1, 0)
		pipe.Set(ctx, lastLeakKey, now, 0)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
