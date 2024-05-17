package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	types "github.com/shariqali-dev/toll-calculator/internal"
)

func main() {
	reciever, err := NewDatareciever()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", reciever.handleWS)
	http.ListenAndServe(":3000", nil)
}

type Datareciever struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

func NewDatareciever() (*Datareciever, error) {
	kafkaTopic := "ubudata"
	p, err := NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	return &Datareciever{
		msgch: make(chan types.OBUData, 128),
		prod:  NewLogMiddleWare(p),
	}, nil
}

func (dr *Datareciever) produceData(data types.OBUData) error {

	return dr.prod.ProduceData(data)
}

func (dr *Datareciever) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsrecieveLoop()
}

func (dr *Datareciever) wsrecieveLoop() {
	fmt.Println("NEW OBU connected,client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
