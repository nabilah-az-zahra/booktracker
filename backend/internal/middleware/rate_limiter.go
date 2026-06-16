package middleware

import (
	"net/http"
	"time"

	redisclient "booktracker/backend/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
		}
		if ip == "" {
			ip = c.ClientIP()
		}

		key := "ratelimit:auth:" + ip
		ctx := c.Request.Context()

		var count int64
		_, err := redisclient.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, time.Minute)
			count = incr.Val()
			return nil
		})

		if err != nil {
			c.Next()
			return
		}

		if count > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests, please try again later", })
			c.Abort()
			return 
		}

		c.Next()
	}
}