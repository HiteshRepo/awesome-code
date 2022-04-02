package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func CreateTopics(topics []string) {
	for _, topic := range topics {
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
		})

		if err != nil {
			log.Fatal(err)
		}

		results, err := adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{{
			Topic: topic,
			NumPartitions: 1,
			ReplicationFactor: 1,
		}})

		if err != nil {
			log.Fatal(err)
		}

		log.Println(results)
	}
}
