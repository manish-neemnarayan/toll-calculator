package main

import (
	"fmt"

	"github.com/manish-neemnarayan/toll-calculator/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
	Invoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("inside aggregate-distance svc")
	i.store.Insert(distance)
	return nil
}

func (i *InvoiceAggregator) Invoice(id int) (inv *types.Invoice, err error) {
	var BaseAmount float64 = 12
	dist, err := i.store.Get(id)
	if err != nil {
		return &types.Invoice{
			OBUId:         id,
			TotalDistance: 0.0,
			TotalAmount:   0,
		}, err
	}

	return &types.Invoice{
		OBUId:         id,
		TotalDistance: dist,
		TotalAmount:   dist * BaseAmount,
	}, nil
}
