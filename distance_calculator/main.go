package main

import (
	"log"
)

const kafkaTopic = "ubudata"

func main() {
	svc := NewCalculatorService()
	KafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	KafkaConsumer.Start()
}
