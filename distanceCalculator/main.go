package main

import "log"

type DistanceCalculator struct {
	consumer DataConsumer
}

func main() {
	d := NewDistanceCalculator()

	d.consumer.consumeData()

}

func NewDistanceCalculator() *DistanceCalculator {
	c, err := NewKafkaConsumer("obutopic")
	if err != nil {
		log.Fatal(err)
	}

	return &DistanceCalculator{
		consumer: c,
	}
}
