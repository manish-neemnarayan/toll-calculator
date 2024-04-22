package main

import (
	"fmt"
	"time"

	"github.com/manish-neemnarayan/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) produceData(data types.OBUData) error {
	fmt.Println("Inside middleware..")
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuId": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("Producing To Kafka")
	}(time.Now())
	return l.next.produceData(data)
}
