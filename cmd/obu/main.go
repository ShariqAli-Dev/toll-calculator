package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

const (
	sendInterval = time.Second * 3
	wsEndpoint   = "ws://127.0.0.1:3000/ws"
)

func genCoord() float64 {
	return rand.Float64() * 100
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(10)
	}
	return ids
}

func main() {
	obuIDS := generateOBUIDS(20)
	connection, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatalf("error init websocket: %v", err)
	}
	for {
		for _, obuId := range obuIDS {
			lat, long := genLatLong()
			data := types.OBUData{OBUID: obuId, Lat: lat, Long: long}
			logrus.WithFields(logrus.Fields{
				"obuID": data.OBUID,
				"lat":   data.Lat,
				"long":  data.Long,
			}).Info("creating obu data")
			if err := connection.WriteJSON(data); err != nil {
				log.Fatalf("error writing ws json: %v", err)
			}
			time.Sleep(sendInterval)
		}
	}
}
