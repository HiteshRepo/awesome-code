package main

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hiteshpattanayak-tw/awesome-code/kafka/consumer"
	"github.com/hiteshpattanayak-tw/awesome-code/kafka/producer"
	"go.uber.org/atomic"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := Start()
	log.Println(err)
}

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	var isTopic1Healthy atomic.Value
	var isTopic2Healthy atomic.Value
	var isTopic3Healthy atomic.Value

	var topic1LastReceived atomic.Value
	var topic2LastReceived atomic.Value
	var topic3LastReceived atomic.Value

	msgChan := make(chan *kafka.Message)

	topics := []string{
		"priority-1",
		"priority-2",
		"priority-3",
	}

	keys := []string {
		"key-1",
		"key-2",
		"key-3",
	}

	messages := []string {
		"msg-1",
		"msg-2",
		"msg-3",
	}

	createTopics(topics, 5*time.Second)

	producers, err := getProducers()
	if err != nil {
		cancel()
		return err
	}

	go startProducing(ctx, producers, topics, keys, messages)

	allConsumer, err := consumer.ProvideKafkaConsumer("failOverGrp", false)
	if err != nil {
		cancel()
		return err
	}

	err = consumer.Subscribe(topics, allConsumer)
	if err != nil {
		cancel()
		return err
	}

	go startProcessing(ctx, msgChan, topics, &isTopic1Healthy, &isTopic2Healthy, &isTopic3Healthy)

	go startConsuming(allConsumer, cancel, topics, &topic1LastReceived, &topic2LastReceived, &topic3LastReceived, &isTopic1Healthy, &isTopic2Healthy, &isTopic3Healthy, msgChan)

	count := 0
	for count < 12 {

		if count == 4 {
			log.Println("closing producer 1")
			producers[0].Close()
		}

		if count == 8 {
			log.Println("closing producer 2")
			producers[1].Close()
		}

		time.Sleep(2 * time.Second)
		count += 1
	}

	<-interrupt()
	cancel()

	log.Println("closing producer 3")
	producers[2].Close()

	err = allConsumer.Close()
	if err != nil {
		return err
	}

	return nil
}

func startProcessing(ctx context.Context, msgChan chan *kafka.Message, topics []string, isTopic1Healthy *atomic.Value, isTopic2Healthy *atomic.Value, isTopic3Healthy *atomic.Value,) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgChan:
			if ok {
				topic := msg.TopicPartition.Topic
				if isTopic1Healthy.Load() != nil && isTopic1Healthy.Load().(bool) {
					if *topic == topics[0] {
						processMessage(msg.Value)
					}
				} else if isTopic2Healthy.Load() != nil && isTopic2Healthy.Load().(bool) {
					if *topic == topics[1] {
						processMessage(msg.Value)
					}
				} else if isTopic3Healthy.Load() != nil && isTopic3Healthy.Load().(bool) {
					if *topic == topics[2] {
						processMessage(msg.Value)
					}
				}
			}
		default:
			continue
		}
	}
}

func startConsuming(allConsumer *kafka.Consumer, cancel context.CancelFunc, topics []string, topic1LastReceived *atomic.Value, topic2LastReceived *atomic.Value, topic3LastReceived *atomic.Value, isTopic1Healthy *atomic.Value, isTopic2Healthy *atomic.Value, isTopic3Healthy *atomic.Value, msgChan chan *kafka.Message) {
	for {
		msg, err := consumer.ReadMessage(allConsumer, "failOverGrp")
		if err != nil {
			cancel()
			return
		}

		topic := msg.TopicPartition.Topic

		switch *topic {
		case topics[0]:
			topic1LastReceived.Store(time.Now())
		case topics[1]:
			topic2LastReceived.Store(time.Now())
		case topics[2]:
			topic3LastReceived.Store(time.Now())
		}

		if topic1LastReceived.Load() != nil {
			if time.Now().Sub(topic1LastReceived.Load().(time.Time)) > (4 * time.Second) {
				isTopic1Healthy.Store(false)
			} else {
				isTopic1Healthy.Store(true)
			}
		}

		if topic2LastReceived.Load() != nil {
			if time.Now().Sub(topic2LastReceived.Load().(time.Time)) > (4 * time.Second) {
				isTopic2Healthy.Store(false)
			} else {
				isTopic2Healthy.Store(true)
			}
		}

		if topic3LastReceived.Load() != nil {
			if time.Now().Sub(topic3LastReceived.Load().(time.Time)) > (4 * time.Second) {
				isTopic3Healthy.Store(false)
			} else {
				isTopic3Healthy.Store(true)
			}
		}

		msgChan <- msg
	}
}

func startProducing(ctx context.Context, producers []*kafka.Producer, topics, keys, messages []string) {
	for {
		err := producer.Produce(topics[0], []byte(keys[0]), []byte(messages[0]), producers[0])
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(1 * time.Second)

		err = producer.Produce(topics[1], []byte(keys[1]), []byte(messages[1]), producers[1])
		if err != nil {
			log.Println(err)
			return
		}


		time.Sleep(1 * time.Second)

		err = producer.Produce(topics[2], []byte(keys[2]), []byte(messages[2]), producers[2])
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(1 * time.Second)

		select {
		case <- ctx.Done():
			return
		default:
			continue
		}
	}
}

func createTopics(topics []string, timeout time.Duration) {
	for _, topic := range topics {
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9094",
		})
		if err != nil {
			log.Println(err)
			return
		}

		results, err := adminClient.CreateTopics(context.Background(),
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

func getProducers() ([]*kafka.Producer, error) {
	producers := make([]*kafka.Producer, 0)

	producer1, err := producer.ProvideKafkaProducer(true)
	if err != nil {
		return nil, err
	}
	producers = append(producers, producer1)

	producer2, err := producer.ProvideKafkaProducer(true)
	if err != nil {
		return nil, err
	}
	producers = append(producers, producer2)

	producer3, err := producer.ProvideKafkaProducer(true)
	if err != nil {
		return nil, err
	}
	producers = append(producers, producer3)
	return producers, nil
}

func processMessage(data []byte) {
	log.Println("Got message : ", string(data))
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}