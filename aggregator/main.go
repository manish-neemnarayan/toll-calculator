package main

import (
	"fmt"
	"net/http"
)

func main() {
	var (
		store = makeStore()
		svc   = NewInvoiceAggregator(store)
	)

	http.HandleFunc("/aggregate", AggregateHandler(svc))
	fmt.Println("agg is working fine...")
}

func AggregateHandler(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
	}
}

func makeStore() Storer {
	return NewMemoryStore()
}
