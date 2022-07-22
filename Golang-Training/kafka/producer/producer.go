package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func ProvideKafkaProducer(isIdempotent bool) (*kafka.Producer, error) {
	kafkaConfig := buildKafkaProducerConfigMap(isIdempotent)
	p, err := kafka.NewProducer(&kafkaConfig)

	if err != nil {
		log.Println("Unable to create kafka producer", err)
		return nil, err
	}

	go deliveryReports(p)

	return p, nil
}

func Produce(topic string, key, message []byte, producer *kafka.Producer) error {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          message,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

func Close(producer *kafka.Producer) error {
	producer.Close()
	return nil
}

func buildKafkaProducerConfigMap(isIdempotent bool) kafka.ConfigMap {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks": -1,
	}

	if isIdempotent {
		_ = kafkaConfig.Set("enable.idempotence=true")
		_ = kafkaConfig.Set("retries=3")
		_ = kafkaConfig.Set("max.in.flight.requests.per.connection=5")
	}

	return kafkaConfig
}

func deliveryReports(producer *kafka.Producer) {
	for e := range producer.Events() {
		if ev, ok := e.(*kafka.Message); ok {
			if ev.TopicPartition.Error != nil {
				log.Printf("Failed to deliver to topic partition: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Successfully delivered to topic partition : %v\n", ev.TopicPartition)
			}
		}
	}
}
