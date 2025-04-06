package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TokenBucketMiddleware(ctx context.Context, rdb *redis.Client, key string, capacity int, refillRate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		lastRefillKey := key + ":last_refill"
		tokensKey := key + ":tokens"

		lastRefill, _ := rdb.Get(ctx, lastRefillKey).Int64()
		tokens, _ := rdb.Get(ctx, tokensKey).Int()

		elapsed := now - lastRefill
		newTokens := int(elapsed) * refillRate
		tokens = min(capacity, tokens+newTokens)

		if tokens <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		pipe := rdb.TxPipeline()
		pipe.Set(ctx, tokensKey, tokens-1, 0)
		pipe.Set(ctx, lastRefillKey, now, 0)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}
