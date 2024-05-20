package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/manish-neemnarayan/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func main() {
	listenAddr := flag.String("port", ":5001", "pass the listen address port value")
	flag.Parse()

	var (
		client     = client.NewHttpClient("http://localhost:3000")
		invHandler = newInvoiceHandler(client)
	)

	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))

	logrus.Infof("gateway is running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func newInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	var obu = r.URL.Query()["obuId"][0]
	obuId, err := strconv.Atoi(obu)
	if err != nil {
		log.Fatal("incorrect obuId")
	}
	invoice, err := h.client.GetInvoice(context.Background(), obuId)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	return writeJson(w, http.StatusOK, invoice)
}

// -----utility method
func writeJson(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error while making request"})
		}
	}
}
