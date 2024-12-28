package queue

import (
    "context"
    "log"
    "fmt"

    "github.com/twmb/franz-go/pkg/kgo"
)

var producer *kgo.Client
var defaultTopic string

func InitProducer(brokers []string) {
    var err error
    producer, err = kgo.NewClient(
        kgo.SeedBrokers(brokers...),
    )
    if err != nil {
        log.Fatalf("Error initializing Kafka producer: %v", err)
    }
}

func InitProducerWithTopic(brokers []string, topic string) {
    InitProducer(brokers)
    defaultTopic = topic
}

func PublishMessage(key, value []byte) error {
    if producer == nil {
        return fmt.Errorf("kafka producer is not initialized")
    }
    if defaultTopic == "" {
        return fmt.Errorf("cannot produce record with no topic and no default topic")
    }
    record := &kgo.Record{
        Topic: defaultTopic,
        Key:   key,
        Value: value,
    }
    return producer.ProduceSync(context.Background(), record).FirstErr()
}

func CloseProducer() {
    if producer != nil {
        producer.Close()
    }
}