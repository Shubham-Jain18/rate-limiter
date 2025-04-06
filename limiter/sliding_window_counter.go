package limiter

import (
	"net/http"
	rlredis "rate-limiter/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SlidingWindowCounterMiddleware(rdb *redis.Client, limit int, windowSizeSec int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = c.ClientIP()
		}

		key := "rl:" + userID

		count, _ := rdb.Get(rlredis.Ctx, key).Int()

		if count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		pipe := rdb.TxPipeline()
		pipe.Incr(rlredis.Ctx, key)
		pipe.Expire(rlredis.Ctx, key, time.Duration(windowSizeSec)*time.Second)
		_, _ = pipe.Exec(rlredis.Ctx)

		c.Next()
	}
}
