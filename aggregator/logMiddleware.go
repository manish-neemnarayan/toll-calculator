package main

import (
	"time"

	"github.com/manish-neemnarayan/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"distance": distance.Value,
			"took":     time.Since(start),
			"obuId":    distance.OBUID,
			"err":      err,
		}).Info("Aggregate Distance: ")
	}(time.Now())

	logrus.Info("Processing and Inserting distance in the storage")
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) Invoice(id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"err":           err,
			"obuId":         inv.OBUId,
			"totalDistance": inv.TotalDistance,
			"totalAmount":   inv.TotalAmount,
		}).Info("Invoice: ")
	}(time.Now())

	inv, err = m.next.Invoice(id)

	return inv, err
}
