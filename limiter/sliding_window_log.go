package limiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SlidingWindowLogMiddleware(ctx context.Context, rdb *redis.Client, key string, limit int, windowSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		windowStart := now - int64(windowSize)

		rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart))
		count, _ := rdb.ZCard(ctx, key).Result()

		if int(count) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		rdb.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: now,
		})
		rdb.Expire(ctx, key, time.Duration(windowSize)*time.Second)

		c.Next()
	}
}
