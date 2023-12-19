package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/shariqali-dev/toll-calc/types"
)

func main() {
	log.Println("the receiver has started")
	receiver := NewDataReceiver()
	http.HandleFunc("/ws", receiver.handleWS)
	http.ListenAndServe(":3000", nil)
	fmt.Println("data receiver working")
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 2048),
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
	fmt.Println("NEW OBU connected,client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Printf("received UBU data from [%d] :: <lat %.2f, long %.2f>\n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data
	}
}
