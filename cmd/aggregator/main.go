package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/shariq/microservice/internal/types"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3001", "the listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		logrus.Error("failed to listen an server server", err)
	}
}

func handleAggregate(_ Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
