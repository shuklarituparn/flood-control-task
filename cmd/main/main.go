package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	config2 "task/internal/config"
	"task/internal/handlers"
	"task/internal/middleware"
)

func main() {

	config, err := config2.LoadConfig("config.yml")
	if err != nil {
		fmt.Println(err)
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
		log.Fatalf("Failed to start server: %v", err)
	}
}
