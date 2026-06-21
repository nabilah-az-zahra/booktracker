package middleware

import (
	"net/http"
	"time"

	redisclient "booktracker/backend/internal/redis"

	"github.com/gin-gonic/gin"
)

func AuthRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip :=  c.ClientIP()
		
		key := "ratelimit:auth:" + ip
		ctx := c.Request.Context()

		pipe := redisclient.Client.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err := pipe.Exec(ctx)

		if err != nil {
			c.Next()
			return
		}

		if incr.Val() > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests, please try again later", })
			c.Abort()
			return 
		}

		c.Next()
	}
}

func GeneralRateLimiter() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("userID")
        if userID == "" {
            c.Next()
            return
        }

        key := "ratelimit:general:" + userID
        ctx := c.Request.Context()

        pipe := redisclient.Client.Pipeline()
        incr := pipe.Incr(ctx, key)
        pipe.Expire(ctx, key, time.Minute)
        _, err := pipe.Exec(ctx)

        if err != nil {
            c.Next()
            return
        }

        if incr.Val() > 100 {
            c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests, please try again later"})
            c.Abort()
            return
        }

        c.Next()
    }
}