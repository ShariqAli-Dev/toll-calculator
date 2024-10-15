package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3001", "the listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewlogMiddleware(NewInvoiceAggregator(store))

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		logrus.Error("failed to listen an server server", err)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
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
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
