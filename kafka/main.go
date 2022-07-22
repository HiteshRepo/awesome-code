package main

import (
	"github.com/hiteshpattanayak-tw/awesome-code/kafka/producer"
	"log"
)

func main() {
	prd, err := producer.ProvideKafkaProducer(false)
	if err != nil {
		log.Println("error creating producer ", err)
		return
	}

	msg := []byte("some message")

	err = producer.Produce("topic.v1", []byte("some key"), msg, prd)
	if err != nil {
		log.Println("error while producing", err)
	}
}
