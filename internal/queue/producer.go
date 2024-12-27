package queue

import (
	"context"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

var producer *kgo.Client

func InitProducer(brokers []string) {
	var err error
	producer, err = kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
}

func PublishMessage(key, value []byte) error {
	err := producer.ProduceSync(context.Background(), &kgo.Record{
		Topic: os.Getenv("KAFKA_TOPIC"),
		Key:   key,
		Value: value,
	}).FirstErr()
	if err != nil {
		log.Printf("Error publishing message to Kafka: %v", err)
		return err
	}
	log.Println("Message successfully published to Kafka")
	return nil
}

func CloseProducer() {
	producer.Close()
}
