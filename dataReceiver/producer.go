package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/manish-neemnarayan/toll-calculator/types"
)

type DataProducer interface {
	produceData(types.OBUData) error
}

type kafkaProducer struct {
	topic    string
	producer *kafka.Producer
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} //else {
				// fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				// }
			}
		}
	}()

	return &kafkaProducer{
		topic:    topic,
		producer: p,
	}, nil

}

func (p *kafkaProducer) produceData(data types.OBUData) error {
	fmt.Println("Inside producer.go producedata...")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)

}
