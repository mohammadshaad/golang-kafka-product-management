package main

import (
	"log"
	"os"

	"github.com/mohammadshaad/zocket/config"
	"github.com/mohammadshaad/zocket/internal/queue"
)

func main() {

	config.LoadConfig()

	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	queue.InitConsumer([]string{brokers}, groupID)

	log.Println("Starting Kafka consumer...")
	queue.ConsumeMessages(queue.ProcessImageMessage)
}
