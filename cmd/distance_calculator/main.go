package main

import (
	"log"

	"github.com/shariqali-dev/toll-calculator/internal/client"
	"github.com/shariqali-dev/toll-calculator/internal/types"
)

func main() {
	svc := NewLogMiddleware(NewCalculatorService())
	aggClient := client.NewClient(types.AGGREGATOR_ENDPOINT)
	kafkaConsumer, err := NewKafkaConsumer(types.KAFKA_TOPIC, svc, aggClient)
	if err != nil {
		log.Fatalf("error init kafka consumer: %v", err)
	}
	kafkaConsumer.Start()

}
