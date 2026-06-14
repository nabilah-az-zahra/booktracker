package middleware

import (
	"context"
	"net/http"
	"time"

	redisclient "booktracker/backend/internal/redis"

	"github.com/gin-gonic/gin"
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
		ctx := context.Background()

		count, err := redisclient.Client.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			redisclient.Client.Expire(ctx, key, time.Minute)
		}

		if count > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests, please try again later", })
			c.Abort()
			return 
		}

		c.Next()
	}
}