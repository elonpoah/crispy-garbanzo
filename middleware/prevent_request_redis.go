package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func PreventRequestsRedis(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-Request-ID"})
			return
		}
		// 使用 Redis SETNX 实现 200ms 防重复
		ok, err := redisClient.SetNX(ctx, requestID, true, 200*time.Millisecond).Result()
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Duplicate request detected"})
			return
		}

		c.Next() // 继续处理请求
	}
}
