package main

import (
	"github.com/hitesh-pattanayak/json_file_processor/stream"
	"log"
)

func main() {
	s := stream.NewJSONStream()
	go func() {
		for data := range s.Watch() {
			if data.Error != nil {
				log.Println(data.Error)
			}
			log.Println(data)
			log.Println("-----------")
		}

	}()
	s.Start("dist/ports.json")
}
