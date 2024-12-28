package main

import (
	"log"
	"os"

	"github.com/mohammadshaad/zocket/config"
	"github.com/mohammadshaad/zocket/internal/queue"
    "github.com/mohammadshaad/zocket/internal/db"
	"github.com/mohammadshaad/zocket/internal/cache"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Database connection
    db.InitDatabase()

	// Initialize Redis
    REDIS_ADDR := os.Getenv("REDIS_ADDR")
    REDIS_USERNAME := os.Getenv("REDIS_USERNAME")
    REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
    cache.InitRedis(REDIS_ADDR, REDIS_USERNAME, REDIS_PASSWORD)

	// Initialize Kafka Consumer
	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	queue.InitConsumer([]string{brokers}, groupID)

	// Initialize S3 Storage
	bucketName := os.Getenv("S3_BUCKET")
	queue.InitS3Storage(bucketName)

	// Start consuming messages
	log.Println("Starting Kafka consumer...")
	queue.ConsumeMessages(queue.ProcessImageMessage)
}
