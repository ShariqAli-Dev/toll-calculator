package main

import (
	"log"

	"github.com/shariqali-dev/toll-calculator/internal/client"
	"github.com/shariqali-dev/toll-calculator/internal/types"
)

func main() {
	svc := NewLogMiddleware(NewCalculatorService())

	httpClient := client.NewHTTPClient(types.AGGREGATOR_ENDPOINT)
	// grpcClient, err := client.NewGRPCClient(types.AGGREGATOR_ENDPOINT)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kafkaConsumer, err := NewKafkaConsumer(types.KAFKA_TOPIC, svc, httpClient)
	if err != nil {
		log.Fatalf("error init kafka consumer: %v", err)
	}
	kafkaConsumer.Start()

}
