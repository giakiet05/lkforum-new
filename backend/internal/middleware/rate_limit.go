package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(redisClient *redis.Client, scope string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		identifier := c.ClientIP()
		if val := strings.TrimSpace(c.GetHeader("X-Forwarded-For")); val != "" {
			identifier = strings.TrimSpace(strings.Split(val, ",")[0])
		}

		ctx, cancel := util.NewDefaultRedisContext()
		defer cancel()

		key := fmt.Sprintf("rate:%s:%s", scope, identifier)
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			_ = redisClient.Expire(ctx, key, window).Err()
		}
		if count > limit {
			dto.SendError(c, http.StatusTooManyRequests, "Too many requests. Please try again later.", "RATE_LIMITED")
			c.Abort()
			return
		}

		c.Next()
	}
}
