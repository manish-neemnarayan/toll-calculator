package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/manish-neemnarayan/toll-calculator/types"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3001", "listen port of http server")
	var (
		store = makeStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHttpTransport(*listenAddr, svc)
}

func makeHttpTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport on port ", listenAddr)

	http.HandleFunc("/aggregate", AggregateHandler(svc))
	http.HandleFunc("/invoice", InvoiceHandler(svc))
	http.ListenAndServe(listenAddr, nil)
}

func InvoiceHandler(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obuId = r.URL.Query()["obuId"][0]

		id, err := strconv.Atoi(obuId)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		inv, err := svc.Invoice(id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		writeJSON(w, http.StatusOK, map[string]any{"data": inv})
	}
}

func AggregateHandler(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeStore() Storer {
	return NewMemoryStore()
}
