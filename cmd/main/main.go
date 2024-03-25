package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	config2 "task/internal/config"
	"task/internal/handlers"
	"task/internal/logger"
	"task/internal/middleware"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	fileLogger := logger.SetupLogger()
	config, err := config2.LoadConfig("config.yml")
	if err != nil {
		fmt.Println(err)
		fileLogger.Printf("error opening config file: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	router := gin.Default()

	router.Use(middleware.RateLimiterMiddleware(client, config))

	router.GET("/persik", handlers.PersikHandler)
	if err := router.Run(":8090"); err != nil {
		fileLogger.Fatalf("Failed to start server: %v", err)
	}
}
