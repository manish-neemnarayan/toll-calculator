package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/manish-neemnarayan/toll-calculator/types"
)

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":9002", nil)
}

type DataReceiver struct {
	conn  *websocket.Conn
	msgch chan types.OBUData
	prod  DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p     DataProducer
		err   error
		topic = "obutopic"
	)
	p, err = NewKafkaProducer(topic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)

	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

// websocket handler
func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("new OBU client connected!")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}

		if err := dr.prod.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
