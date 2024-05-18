package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	types "github.com/shariqali-dev/toll-calculator/internal"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "thi listen addresss of the HTTP serve")
	flag.Parse()
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/post", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
