package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type DataReciever struct {
	msgChan chan types.OBUData
	conn    *websocket.Conn
	prod    DataProducer
}

func NewDataReciever() (*DataReciever, error) {
	var (
		p   DataProducer
		err error
	)

	p, err = NewKafkaProducer(types.KAFKA_TOPIC)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)

	return &DataReciever{
		msgChan: make(chan types.OBUData, 128),
		prod:    p,
	}, nil
}

func (dr *DataReciever) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReciever) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1028, 1028)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsRecieveLoop()
}

func (dr *DataReciever) wsRecieveLoop() {
	fmt.Println("NEW OBU CONNECTED!")

	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Printf("ws recieve read json error: %v", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Printf("kafka produce error: %v\n", err)
		}
	}
}

func main() {
	reciever, err := NewDataReciever()
	if err != nil {
		log.Fatalf("error creating data reciever: %v\n", err)
	}

	http.HandleFunc("/ws", reciever.handleWS)
	http.ListenAndServe(":3000", nil)
}
