package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	types "github.com/shariqali-dev/toll-calculator/internal"
	"github.com/sirupsen/logrus"
)

// this can also be called kafka transport
type KafkaConsumer struct {
	consumer     *kafka.Consumer
	isRunning    bool
	calcServicer CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}
	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:     c,
		calcServicer: svc,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}
		distance, err := c.calcServicer.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error serialization error: %s", err)
			continue
		}
		// fmt.Printf("distance %.2f\n", distance)
		_ = distance
	}
}
