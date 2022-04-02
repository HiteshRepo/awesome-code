package main

import (
	"github.com/hiteshrepo/awesome-code/kowl/kafka"
	user_v1 "github.com/hiteshrepo/awesome-code/kowl/proto"
	"log"
	"sync"
)

func main() {
	kafka.CreateTopics([]string{"user.v1"})
	producer := kafka.ProvideKafkaProducer()

	user := user_v1.User{
		Name: "Hitesh",
		City: "Pune",
		Age:  28,
	}

	userMsg, err := user.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	producer.Produce("user.v1", []byte(user.Name), userMsg)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
