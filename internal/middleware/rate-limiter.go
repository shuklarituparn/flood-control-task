package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"task/internal/config"
	"task/internal/logger"
	redis_impelementation "task/internal/redis"
)

func RateLimiterMiddleware(client *redis.Client, Config *config.RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileLogger := logger.SetupLogger()
		userIDfromUrl := c.Query("userID")
		userTypefromUrl := c.Query("userType")
		fmt.Println(userIDfromUrl)
		userID, _ := strconv.ParseInt(userIDfromUrl, 10, 64)

		if userIDfromUrl == "" && userTypefromUrl == "" {
			DefaultSlidingWindow := redis_impelementation.NewSlidingWindow(client, "rl:default", 23, float64(Config.RateLimiter.DefaultRateLimit.Rate), Config.RateLimiter.DefaultRateLimit.WindowSeconds)
			check, err := DefaultSlidingWindow.Check(context.Background(), 23)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				fileLogger.Printf("server error: %v", err.Error())
				return
			}
			if !check {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
			c.Next()
		}
		if userIDfromUrl != "" && userTypefromUrl == "" {
			DefaultSlidingWindow := redis_impelementation.NewSlidingWindow(client, "rl:default", userID, float64(Config.RateLimiter.DefaultRateLimit.Rate), Config.RateLimiter.DefaultRateLimit.WindowSeconds)
			check, err := DefaultSlidingWindow.Check(context.Background(), userID)
			if err != nil {
				fileLogger.Printf("server error: %v", err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !check {

				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
			c.Next()
		}
		if userIDfromUrl == "" && userTypefromUrl != "" {
			NormalSlidingWindow := redis_impelementation.NewSlidingWindow(client, Config.RateLimiter.UserTypes["normal"].KeyPrefix, userID, float64(Config.RateLimiter.UserTypes["normal"].RateLimit.Rate), Config.RateLimiter.UserTypes["normal"].RateLimit.WindowSeconds)
			check, err := NormalSlidingWindow.Check(context.Background(), userID)
			if err != nil {
				fileLogger.Printf("server error: %v", err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !check {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
			c.Next()
		}
		if userIDfromUrl != "" && userTypefromUrl != "" {
			PremiumSlidingWindow := redis_impelementation.NewSlidingWindow(client, Config.RateLimiter.UserTypes[userTypefromUrl].KeyPrefix, userID, float64(Config.RateLimiter.UserTypes[userTypefromUrl].RateLimit.Rate), Config.RateLimiter.UserTypes[userTypefromUrl].RateLimit.WindowSeconds)
			check, err := PremiumSlidingWindow.Check(context.Background(), userID)
			if err != nil {
				fileLogger.Printf("server error: %v", err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !check {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				return
			}
			c.Next()
		}

	}
}
