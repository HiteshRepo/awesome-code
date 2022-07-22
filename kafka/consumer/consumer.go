package consumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

func ProvideKafkaConsumer(groupName string, disableAutoCommit bool) (*kafka.Consumer, error) {
	kafkaConfig := buildKafkaConsumerConfigMap(groupName, disableAutoCommit)
	c, err := kafka.NewConsumer(&kafkaConfig)

	if err != nil {
		log.Println("Unable to create kafka consumer", err)
		return nil, err
	}

	return c, nil
}

func Subscribe(topics []string, consumer *kafka.Consumer) error {
	err := consumer.SubscribeTopics(topics, nil)

	if err != nil {
		log.Println("Unable to Subscribe to topics", err)
		return err
	}
	log.Println("Subscribed successfully")
	return nil
}

func ReadMessage(consumer *kafka.Consumer, groupName string) (*kafka.Message, error) {
	msg, err := consumer.ReadMessage(5 * time.Second)

	if err == nil {
		log.Println(fmt.Sprintf("Message on %s: %s of group: %s\n", msg.TopicPartition, string(msg.Value), groupName))
		return msg, nil
	}

	log.Println("Unable to read message", err)
	return nil, err
}

func Close(consumer *kafka.Consumer) error {
	err := consumer.Close()
	if err != nil {
		return err
	}
	return nil
}

func buildKafkaConsumerConfigMap(groupName string, disableAutoCommit bool) kafka.ConfigMap {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9094",
		"group.id":                 groupName,
		"auto.offset.reset":        "earliest",
	}

	if disableAutoCommit {
		_ = kafkaConfig.Set("enable.auto.commit=false")
		_ = kafkaConfig.Set("auto.commit.interval.ms=0")
		_ = kafkaConfig.Set("enable.auto.offset.store=false")
	}

	return kafkaConfig
}
