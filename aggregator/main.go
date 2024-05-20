package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/manish-neemnarayan/toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "listen port of HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "listen port of GRPC server")
	var (
		store = makeStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc))
	}()
	// time.Sleep(time.Second * 5)
	// c, err := client.NewGRPCClient(*grpcListenAddr)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.Aggregate(context.Background(), &types.AggregateRequest{
	// 	OBUID: 1,
	// 	Value: 58.23,
	// 	Unix:  time.Now().Unix(),
	// })

	makeHttpTransport(*httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport on port ", listenAddr)
	// Make a TCP listener
	fmt.Println(listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("in tcp error")
		return err
	}

	defer func() {
		fmt.Println("stopping GRPC transport")
		ln.Close()
	}()

	//Make a new grpc server
	server := grpc.NewServer([]grpc.ServerOption{}...)

	// Register grpc sercer implementation to the grpc package.
	types.RegisterAggregatorServer(server, NewGRPCServer(svc))
	return server.Serve(ln)
}

func makeHttpTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport on port ", listenAddr)

	http.HandleFunc("/aggregate", AggregateHandler(svc))
	http.HandleFunc("/invoice", InvoiceHandler(svc))

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func InvoiceHandler(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obuId = r.URL.Query()["obuId"][0]
		fmt.Println(obuId)
		id, err := strconv.Atoi(obuId)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		inv, err := svc.Invoice(id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		writeJSON(w, http.StatusOK, inv)
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
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeStore() Storer {
	return NewMemoryStore()
}
