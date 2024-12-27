package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohammadshaad/zocket/config"
	"github.com/mohammadshaad/zocket/internal/api"
	"github.com/mohammadshaad/zocket/internal/db"
	"github.com/mohammadshaad/zocket/internal/queue"
	"github.com/mohammadshaad/zocket/internal/cache"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Database
	db.InitDatabase()
	db.Migrate()

	// Initialize Kafka Producer
	brokers := os.Getenv("KAFKA_BROKERS")
	queue.InitProducer([]string{brokers})
	defer queue.CloseProducer()

	// Initialize Redis (optional)
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	USERNAME := os.Getenv("REDIS_USERNAME")
	cache.InitRedis(REDIS_ADDR, USERNAME, REDIS_PASSWORD)
 
	router := gin.Default()

	api.SetupRoutes(router)

	log.Println("API server running on port 8080")
	log.Fatal(router.Run(":8080"))
}
