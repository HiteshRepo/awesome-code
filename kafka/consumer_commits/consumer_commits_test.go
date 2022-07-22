package consumer_commits_test

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

type consumerCommitsSuite struct {
	suite.Suite

	ctx context.Context
	topic string

	producer *kafka.Producer
}

func TestConsumerCommitSuite(t *testing.T) {
	suite.Run(t, new(consumerCommitsSuite))
}

func (suite *consumerCommitsSuite) SetupTest() {
	var err error
	suite.ctx = context.Background()
	suite.topic = "test-topic"
	const MaxUint = ^uint(0)
	const MaxInt = int(MaxUint >> 1)

	suite.createTopics([]string{suite.topic}, 10*time.Second)

	suite.producer, err = producer.ProvideKafkaProducer(false)
	suite.Require().NoError(err)
}

func (suite *consumerCommitsSuite) TearDownTest() {
	if suite.producer != nil {
		err := producer.Close(suite.producer)
		suite.Require().NoError(err)
	}
}

func (suite *consumerCommitsSuite) TestEnabledAutoCommit() {
	var err error
	consumer1, err := consumer.ProvideKafkaConsumer("group-1", false)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().EqualError(err, "Local: Timed out")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)
}

func (suite *consumerCommitsSuite) TestDisabledAutoCommitViaMessage() {
	var err error
	consumer1, err := consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	log.Println("skipping commit")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)

	// Restart to read the uncommitted message
	consumer1, err = consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	msg, err := consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	log.Println("committing")

	_, err = consumer1.CommitMessage(msg)
	suite.Require().NoError(err)

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)

	// Restart but this time no message left uncommitted
	consumer1, err = consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().EqualError(err, "Local: Timed out")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)
}

func (suite *consumerCommitsSuite) TestDisabledAutoCommitViaOffset() {
	var err error
	consumer1, err := consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	log.Println("skipping commit")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)

	// Restart to read the uncommitted message
	consumer1, err = consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	msg, err := consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	log.Println("committing")

	offsets := []kafka.TopicPartition{msg.TopicPartition}
	offsets[0].Offset++

	_, err = consumer1.CommitOffsets(offsets)

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)

	// Restart but this time no message left uncommitted
	consumer1, err = consumer.ProvideKafkaConsumer("group-1", true)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().EqualError(err, "Local: Timed out")

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)
}

func (suite *consumerCommitsSuite) TestConsumersInDifferentGroups() {
	var err error
	consumer1, err := consumer.ProvideKafkaConsumer("group-1", false)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer1)
	suite.Require().NoError(err)

	suite.flushMessages(consumer1, "group-1")

	consumer2, err := consumer.ProvideKafkaConsumer("group-2", false)
	suite.Require().NoError(err)

	err = consumer.Subscribe([]string{suite.topic}, consumer2)
	suite.Require().NoError(err)

	suite.flushMessages(consumer2, "group-2")

	times := 1
	suite.produceMessageForGivenTimes("Some key", "some message", times)

	_, err = consumer.ReadMessage(consumer1, "group-1")
	suite.Require().NoError(err)

	_, err = consumer.ReadMessage(consumer2, "group-2")
	suite.Require().NoError(err)

	err = consumer.Close(consumer1)
	suite.Require().NoError(err)

	err = consumer.Close(consumer2)
	suite.Require().NoError(err)
}

func (suite *consumerCommitsSuite) flushMessages(suiteConsumer *kafka.Consumer, groupName string) {
	log.Println("------------flushing starts--------------")
	for {
		_, err := consumer.ReadMessage(suiteConsumer, groupName)
		if err != nil {
			log.Println("------------flushing over--------------")
			return
		}
	}
}

func (suite *consumerCommitsSuite) createTopics(topics []string, timeout time.Duration) {
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

func (suite *consumerCommitsSuite) produceMessageForGivenTimes(key, message string ,times int) {
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


