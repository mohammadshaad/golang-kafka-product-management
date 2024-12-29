#!/bin/sh

# Wait for Kafka to be ready
while ! nc -z kafka 9092; do
  echo "Waiting for Kafka..."
  sleep 2
done

# Create the Kafka topic
kafka-topics.sh --create --topic image_processing --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1

echo "Kafka topic 'image_processing' created."