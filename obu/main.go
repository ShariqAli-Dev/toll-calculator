package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const sendInterval = 1 * time.Second

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	obuIDS := generateOBUIDS(20)
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()
			data := OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Println("%+v\n", data)
			fmt.Println(genCoord())
		}
		fmt.Println(genCoord())
		time.Sleep(sendInterval)
	}
}

func init() {
	// rand.Seed(time.Now().UnixNano())
}
