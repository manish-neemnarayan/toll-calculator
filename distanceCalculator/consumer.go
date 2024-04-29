package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/manish-neemnarayan/toll-calculator/aggregator/client"
	"github.com/manish-neemnarayan/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type DataConsumer interface {
	consumeData()
}

type KafkaConsumer struct {
	topic       string
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string) (DataConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatal(err)
	}

	c.SubscribeTopics([]string{topic}, nil)

	calcSVC := NewCalculatorService()
	nextCalcSVC := NewLogMiddleware(calcSVC)
	aggClient := client.NewClient("http://localhost:3001/aggregate")
	// A signal handler or similar could be used to set this to false to break the loop.
	return &KafkaConsumer{
		consumer:    c,
		topic:       topic,
		calcService: nextCalcSVC,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) consumeData() {
	c.isRunning = true

	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error: %s", err)
		}

		var data types.OBUData

		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			logrus.WithFields(logrus.Fields{
				"err":       err,
				"requestId": data.OBUID,
			})

			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}

		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}

		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error: %v", err)
			continue
		}

	}
}
