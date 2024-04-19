package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/manish-neemnarayan/toll-calculator/types"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for i := 0; i < 15; i++ {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("testing producer..."),
		}, nil)
	}
	// return
	// p.Flush(1000 * 15)
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":9002", nil)
}

type DataReceiver struct {
	conn  *websocket.Conn
	msgch chan types.OBUData
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
	}
}

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

		fmt.Printf("%+v\n", data)
		// dr.msgch <- data //when channel is full our code is blocked here
	}
}
