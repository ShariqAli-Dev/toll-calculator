package main

import (
	"log"
)

const kafkaTopic = "ubudata"

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewlogMiddleware(svc)

	KafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	KafkaConsumer.Start()
}
