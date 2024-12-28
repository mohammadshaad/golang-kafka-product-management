package queue

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/twmb/franz-go/pkg/kgo"
)

var consumer *kgo.Client

func InitConsumer(brokers []string, groupID string) {
    var err error
    consumer, err = kgo.NewClient(
        kgo.SeedBrokers(brokers...),
        kgo.ConsumerGroup(groupID),
        kgo.ConsumeTopics(os.Getenv("KAFKA_TOPIC")),
    )
    if err != nil {
        log.Fatalf("Error initializing Kafka consumer: %v", err)
    }
}

func ConsumeMessages(processFunc func(key, value []byte) error) {
    defer consumer.Close()

    for {
        fetches := consumer.PollFetches(context.Background())

        if fetchErrs := fetches.Errors(); len(fetchErrs) > 0 {
            for _, err := range fetchErrs {
                log.Printf("Error consuming Kafka message: %v", err)
            }
            continue
        }

        iter := fetches.RecordIter()
        for !iter.Done() {
            record := iter.Next()
            log.Printf("Received message: Key=%s, Value=%s", record.Key, record.Value)

            if err := processFunc(record.Key, record.Value); err != nil {
                log.Printf("Error processing Kafka message: %v", err)
            }
        }
        time.Sleep(time.Second)
    }
}