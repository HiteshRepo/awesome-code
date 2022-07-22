package main

import (
	"encoding/json"
	"github.com/hiteshpattanayak-tw/golang-training/kafka/consumer"
	"github.com/hiteshpattanayak-tw/golang-training/kafka/producer"
	"log"
	"sync"
	"time"
)

func main() {
	go ConsumeCarMessages()
	ProducerCarMessages()

	//go ConsumeIceCreamMessages()
	//ProducerIceCreamMessages()
}

type Car struct {
	Manufacturer string
	Type         string
}

func ProducerCarMessages() {
	prd, err := producer.ProvideKafkaProducer(false)
	if err != nil {
		log.Println("error creating producer ", err)
		return
	}

	carMessage1 := Car{
		Manufacturer: "TATA",
		Type:         "Diesel",
	}

	carMessage2 := Car{
		Manufacturer: "TESLA",
		Type:         "Electric",
	}

	carMessages := []Car{carMessage1, carMessage2}

	for _, cm := range carMessages {
		b, _ := json.Marshal(cm)
		err = producer.Produce("cars.v1", []byte(cm.Manufacturer), b, prd)
		if err != nil {
			log.Println("error while producing car message ", err)
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func ConsumeCarMessages() {
	c, err := consumer.ProvideKafkaConsumer("cars", false)
	if err != nil {
		log.Fatal("error creating car consumer ", err)
	}

	err = c.Subscribe("cars.v1", nil)
	if err != nil {
		log.Fatal("error subscribing car topic ", err)
	}

	for {
		msg, err := c.ReadMessage(30 * time.Second)
		if err != nil {
			log.Println("failed to read message", err)
		} else {
			log.Println("Car consumer: Message Received : ", string(msg.Value))
		}
	}
}

type IceCream struct {
	Brand  string
	Flavor string
}

func ProducerIceCreamMessages() {
	prd, err := producer.ProvideKafkaProducer(false)
	if err != nil {
		log.Println("error creating producer ", err)
		return
	}

	icecream1 := IceCream{
		Brand:  "Kwality Walls",
		Flavor: "Chocolate",
	}

	icecream2 := IceCream{
		Brand:  "Baskin Robins",
		Flavor: "Vanilla",
	}

	iceCreamMessages := []IceCream{icecream1, icecream2}

	for _, im := range iceCreamMessages {
		b, _ := json.Marshal(im)
		err = producer.Produce("icecream.v1", []byte(im.Brand), b, prd)
		if err != nil {
			log.Println("error while producing ice cream message ", err)
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func ConsumeIceCreamMessages() {
	c, err := consumer.ProvideKafkaConsumer("icecream", false)
	if err != nil {
		log.Fatal("error creating ice cream consumer ", err)
	}

	err = c.Subscribe("icecream.v1", nil)
	if err != nil {
		log.Fatal("error subscribing ice cream topic ", err)
	}

	for {
		msg, err := c.ReadMessage(30 * time.Second)
		if err != nil {
			log.Println("failed to read message", err)
		} else {
			log.Println("Ice Cream consumer: Message Received : ", string(msg.Value))
		}
	}
}
