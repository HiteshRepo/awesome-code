package suites

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	kafkaTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	KafkaBrokerPort = "9092"
	KafkaClientPort = "9093"
	ZookeeperPort = "2181"
	ZookeeperTickTime = "2000"
	ZookeeperImage = "confluent/cp-zookeeper:5.4.6"
	KafkaImage = "confluent/cp-kafka:5.4.6"
	WaitImage = "waisbrot/wait"
	KafkaStartFileMode = 0700
	TopicCreationTimeoutInSeconds = 10
	KafkaReadTimeout = 30
)

type KafkaSuite struct {
	suite.Suite

	ctx     context.Context
	network *testcontainers.DockerNetwork
	rootDir string

	kafkaContainer     testcontainers.Container
	zookeeperContainer testcontainers.Container
	waitContainer      testcontainers.Container
	producer           *kafka.Producer
}

func (suite *KafkaSuite) GetContainerHost() string {
	host, err := suite.kafkaContainer.Host(suite.ctx)
	suite.Require().NoError(err)

	return host
}

func (suite *KafkaSuite) GetContainerMappedPort() nat.Port {
	port, err := suite.kafkaContainer.MappedPort(suite.ctx, KafkaClientPort)
	suite.Require().NoError(err)

	return port
}

func (suite *KafkaSuite) SetNetwork(network *testcontainers.DockerNetwork) {
	suite.network = network
}

func (suite *KafkaSuite) SetCtx(ctx context.Context) {
	suite.ctx = ctx
}

func (suite *KafkaSuite) SetRootDir(rootDir string) {
	suite.rootDir = rootDir
}

func (suite *KafkaSuite) SetupSuite() {
	suite.createZookeeperContainer()
	suite.createKafkaContainer()
	suite.createWaitContainer()

	err := suite.zookeeperContainer.Start(suite.ctx)
	suite.Require().NoError(err)

	err = suite.kafkaContainer.Start(suite.ctx)
	suite.Require().NoError(err)

	suite.startKafka()

	err = suite.waitContainer.Start(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *KafkaSuite) TearDownSuite() {
	suite.terminateContainer(suite.waitContainer)
	suite.terminateContainer(suite.kafkaContainer)
	suite.terminateContainer(suite.zookeeperContainer)
}

func (suite *KafkaSuite) CreateTopics(topics []string) {
	for _, topic := range topics {
		adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
			"bootstrap.servers": suite.getKafkaHost(),
		})
		suite.Require().NoError(err)

		results, err := adminClient.CreateTopics(suite.ctx, []kafka.TopicSpecification{{
			Topic: topic,
			NumPartitions: 1,
			ReplicationFactor: 1,
		}},
		kafka.SetAdminOperationTimeout(TopicCreationTimeoutInSeconds * time.Second))

		suite.Require().NoError(err)

		for _, res := range results {
			suite.Require().False(res.Error.Code() != kafka.ErrNoError && res.Error.Code() != kafka.ErrTopicAlreadyExists, "topic creation failed")
		}
	}
}

func (suite *KafkaSuite) ConsumeMessage(topic string) *kafka.Message {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": suite.getKafkaHost(),
		"group.id": "myGroup",
		"auto.offset.reset": "earliest",
	})
	suite.Require().NoError(err)

	defer func() {
		err = c.Close()
		suite.Require().NoError(err)
	}()

	err = c.SubscribeTopics([]string{topic}, nil)
	suite.Require().NoError(err)

	msg, err := c.ReadMessage(KafkaReadTimeout * time.Second)
	suite.Require().NoError(err)
	return msg
}

func (suite *KafkaSuite) ProduceMessage(topic string, key, msg []byte) uint64 {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": suite.getKafkaHost(),
	})
	suite.Require().NoError(err)
	suite.producer = p

	go suite.deliveryReports()

	defer func() {
		p.Close()
	}()

	parentSpan := tracer.StartSpan("test_parent_span")

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
		Key: key,
	}
	carrier := kafkaTrace.NewMessageCarrier(message)
	err = tracer.Inject(parentSpan.Context(), carrier)
	suite.Require().NoError(err)

	id := parentSpan.Context().TraceID()

	err = p.Produce(message, nil)
	suite.Require().NoError(err)

	p.Flush(100)

	return id
}

func (suite *KafkaSuite) deliveryReports() {
	for e := range suite.producer.Events() {
		if ev, ok := e.(*kafka.Message); ok {
			if ev.TopicPartition.Error != nil {
				log.Println(fmt.Sprintf("Failed to deliver to topic %v", ev.TopicPartition))
			} else {
				log.Println(fmt.Sprintf("Successfully delivered to topic %v", ev.TopicPartition))
			}
		}
	}
}

func (suite *KafkaSuite) terminateContainer(container testcontainers.Container) {
	if container != nil {
		err := container.Terminate(suite.ctx)
		suite.Require().NoError(err)
	}
}

func (suite *KafkaSuite) createZookeeperContainer() {
	req := testcontainers.ContainerRequest{
		Image:        ZookeeperImage,
		ExposedPorts: []string{ZookeeperPort},
		Env: map[string]string{"ZOOKEEPER_CLIENT_PORT": ZookeeperPort, "ZOOKEEPER_TICK_TIME": ZookeeperTickTime},
		Networks: []string{suite.network.Name},
		NetworkAliases: map[string][]string{suite.network.Name: {"zookeeper"}},
	}
	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	suite.Require().NoError(err)

	suite.zookeeperContainer = container
}

func (suite *KafkaSuite) createKafkaContainer() {
	req := testcontainers.ContainerRequest{
		Image:        KafkaImage,
		ExposedPorts: []string{KafkaClientPort},
		Env: map[string]string {
			"KAFKA_BROKER_ID": "1",
			"KAFKA_ZOOKEEPER_CONNECT": "zookeeper:" + ZookeeperPort,
			"KAFKA_LISTENERS": "PLAINTEXT://0.0.0.0:" + KafkaClientPort + ",BROKER://0.0.0.0:" + KafkaBrokerPort,
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME": "BROKER",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
		},
		Networks: []string{suite.network.Name},
		NetworkAliases: map[string][]string{suite.network.Name: {"kafka"}},
		Cmd: []string{"sh", "-c", "while [ ! -f /testcontainers_start.sh ]; do sleep 0.1; done; /testcontainers_start.sh"},
	}
	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	suite.Require().NoError(err)

	suite.kafkaContainer = container
}

func (suite *KafkaSuite) createWaitContainer() {
	req := testcontainers.ContainerRequest{
		Image:        WaitImage,
		Env: map[string]string{"TARGETS": "kafka:" + KafkaBrokerPort, "TIMEOUT": "120"},
		Networks: []string{suite.network.Name},
		NetworkAliases: map[string][]string{suite.network.Name: {"wait"}},
		WaitingFor: wait.ForLog("Everything is up").WithStartupTimeout(2 * time.Minute),
	}
	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	suite.Require().NoError(err)

	suite.waitContainer = container
}

func (suite *KafkaSuite) getKafkaHost() string {
	return suite.GetContainerHost() + ":" + suite.GetContainerMappedPort().Port()
}

func (suite *KafkaSuite) startKafka() {
	kafkaStartFile, err := ioutil.TempFile("", "testcontainers_start.sh")
	suite.Require().NoError(err)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(kafkaStartFile.Name())

	exposedHost := suite.getKafkaHost()

	tmpl := `
				#!/bin/bash
				export KAFKA_ADVERTISED_LISTENERS='PLAINTEXT://%s,,BROKER://kafka:%s'
				. /etc/confluent/docker/bash-config
				/etc/confluent/docker/configure
				/etc/confluent/docker/launch
			`

	_, err = kafkaStartFile.WriteString(fmt.Sprintf(tmpl, exposedHost, KafkaBrokerPort))
	suite.Require().NoError(err)

	err = suite.kafkaContainer.CopyFileToContainer(suite.ctx, kafkaStartFile.Name(), "testcontainers_start.sh", KafkaStartFileMode)
	suite.Require().NoError(err)
}

