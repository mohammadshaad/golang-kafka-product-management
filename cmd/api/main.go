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

	// Print environment variables for debugging
	printEnvVariables()

	// Initialize Database
	db.InitDatabase()
	db.Migrate()

	// Initialize Kafka Producer
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC environment variable is not set")
	}
	queue.InitProducerWithTopic([]string{brokers}, topic)
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

func printEnvVariables() {
	log.Printf("DATABASE_DSN: %s", os.Getenv("DATABASE_DSN"))
	log.Printf("KAFKA_BROKERS: %s", os.Getenv("KAFKA_BROKERS"))
	log.Printf("KAFKA_TOPIC: %s", os.Getenv("KAFKA_TOPIC"))
	log.Printf("REDIS_ADDR: %s", os.Getenv("REDIS_ADDR"))
	log.Printf("REDIS_PASSWORD: %s", os.Getenv("REDIS_PASSWORD"))
	log.Printf("REDIS_USERNAME: %s", os.Getenv("REDIS_USERNAME"))
	log.Printf("AWS_REGION: %s", os.Getenv("AWS_REGION"))
	log.Printf("S3_BUCKET: %s", os.Getenv("S3_BUCKET"))
	log.Printf("AWS_ACCESS_KEY_ID: %s", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Printf("AWS_SECRET_ACCESS_KEY: %s", os.Getenv("AWS_SECRET_ACCESS_KEY"))
}
