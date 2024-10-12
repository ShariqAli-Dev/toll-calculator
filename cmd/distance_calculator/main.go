package main

import (
	"log"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

func main() {
	svc := NewLogMiddleware(NewCalculatorService())
	kafkaConsumer, err := NewKafkaConsumer(types.KAFKA_TOPIC, svc)
	if err != nil {
		log.Fatalf("error init kafka consumer: %v", err)
	}
	kafkaConsumer.Start()

}
