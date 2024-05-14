package main

import (
	"context"
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
	aggClient   client.Client
}

func NewKafkaConsumer(topic string) (DataConsumer, error) {
	const aggregatorEndpoint = ":3001"

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
	// aggClient := client.NewHttpClient("http://localhost:3001/aggregate")

	grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	// A signal handler or similar could be used to set this to false to break the loop.
	return &KafkaConsumer{
		consumer:    c,
		topic:       topic,
		calcService: nextCalcSVC,
		aggClient:   grpcClient,
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

		req := types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: int32(data.OBUID),
		}

		if err := c.aggClient.Aggregate(context.Background(), &req); err != nil {
			logrus.Errorf("aggregate error: %v", err)
			continue
		}

	}
}
