package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Producer struct {
	producer *kafka.Producer
}

func ProvideKafkaProducer() *Producer {
	kafkaConfig := buildKafkaProducerConfigMap()

	p, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}

	go deliveryReports(p)

	return &Producer{producer: p}
}

func (p *Producer) Produce(topic string, key, message []byte) {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Key:            key,
	}, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func (p *Producer) Close() {
	p.producer.Close()
}

func buildKafkaProducerConfigMap() *kafka.ConfigMap {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":                     "localhost:9092",
		"acks":                                  -1,
		"enable.idempotence":                    true,
		"retries":                               3,
		"max.in.flight.requests.per.connection": 5,
	}

	return kafkaConfig
}

func deliveryReports(producer *kafka.Producer) {
	for e := range producer.Events() {
		if ev, ok := e.(*kafka.Message); ok {
			if ev.TopicPartition.Error != nil {
				log.Printf("Failed to deliver to topic partition: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Successfully delivered to topic partition: %v\n", ev.TopicPartition)
			}
		}
	}
}
