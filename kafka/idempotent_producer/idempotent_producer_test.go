package idempotent_producer_test

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hiteshpattanayak-tw/awesome-code/kafka/consumer"
	"github.com/hiteshpattanayak-tw/awesome-code/kafka/producer"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type idempotentProducerTest struct {
	suite.Suite

	ctx context.Context
	topic string

	producer *kafka.Producer
}

func TestIdempotentProducerSuite(t *testing.T) {
	suite.Run(t, new(idempotentProducerTest))
}

func (suite *idempotentProducerTest) SetupTest() {
	suite.ctx = context.Background()
	suite.topic = "test-topic"
	const MaxUint = ^uint(0)
	const MaxInt = int(MaxUint >> 1)

	suite.createTopics([]string{suite.topic}, 10*time.Second)
}

func (suite *idempotentProducerTest) TearDownTest() {
	if suite.producer != nil {
		err := producer.Close(suite.producer)
		suite.Require().NoError(err)
	}
}

func (suite *idempotentProducerTest) TestMessageDeliveredMultipleTimeAsIdempotencyDisabled() {
	var err error
	suite.producer, err = producer.ProvideKafkaProducer(false)
	suite.Require().NoError(err)

	consumer1, err := consumer.ProvideKafkaConsumer("group-1", false)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)
	//suite.Require().EqualError(err, "Local: Timed out")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)
}

func (suite *idempotentProducerTest) TestMessageDeliveredSingleTimeAsIdempotencyEnabled() {
	var err error
	suite.producer, err = producer.ProvideKafkaProducer(true)
	suite.Require().NoError(err)

	consumer1, err := consumer.ProvideKafkaConsumer("group-1", false)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	err = producer.Close(suite.producer)
	suite.Require().NoError(err)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	suite.producer, err = producer.ProvideKafkaProducer(true)
	suite.Require().NoError(err)

	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().EqualError(err, "Local: Timed out")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)
}

func (suite *idempotentProducerTest) flushMessages(suiteConsumer *kafka.Consumer, groupName string) {
	log.Println("------------flushing starts--------------")
	for {
		_, err := consumer.ReadMessage(suiteConsumer, groupName)
		if err != nil {
			log.Println("------------flushing over--------------")
			return
		}
	}
}

func (suite *idempotentProducerTest) createTopics(topics []string, timeout time.Duration) {
	for _, topic := range topics {
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9094",
		})
		if err != nil {
			log.Println(err)
			return
		}

		results, err := adminClient.CreateTopics(suite.ctx,
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1}},
			kafka.SetAdminOperationTimeout(timeout))

		if err != nil {
			log.Println(err)
			return
		}

		for _, result := range results {
			if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
				log.Println("topic creation failed")
			}
		}
	}
}

func (suite *idempotentProducerTest) produceMessageForGivenTimes(key, message string ,times int) {
	count := 0
	for count < times {
		time.Sleep(2 * time.Second)
		message := fmt.Sprintf("%s %d", message, count)
		err := producer.Produce(suite.topic, []byte(key), []byte(message), suite.producer)
		if err != nil {
			log.Println(err)
			return
		}
		count += 1
	}
}