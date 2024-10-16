package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/shariqali-dev/toll-calculator/internal/client"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.HTTPClient
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.HTTPClient) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		c.Close()
		return nil, err
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
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
		if err == nil {
			// fmt.Printf("mesage on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			logrus.Errorf("kafa consumer error: %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("json serialization error: %s", err)
			continue
		}
		distance := c.calcService.CalculateDistance(data)
		req := types.Distance{
			Values: distance,
			Unix:   time.Now().UnixNano(),
			OBUID:  data.OBUID,
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error: %s", err)
		}
	}

}
