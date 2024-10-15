package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

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
	mux := http.NewServeMux()
	mux.Handle("POST /aggregate", handleAggregate(svc))
	mux.Handle("GET /invoice", handleGetInvoice(svc))

	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		logrus.Error("failed to listen an server server\n", err)
		log.Fatal(err)
	}
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID, err := strconv.Atoi(r.URL.Query().Get("obuID"))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		invoice, err := service.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusInternalServerError, invoice)
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
		writeJSON(w, http.StatusOK, distance)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
