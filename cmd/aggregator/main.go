package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var (
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		store          = makeStore(os.Getenv("AGG_STORE_TYPE"))
		svc            = NewlogMiddleware(NewMetricsMiddleware(NewInvoiceAggregator(store)))
	)

	if httpListenAddr == "" || grpcListenAddr == "" {
		logrus.WithFields(logrus.Fields{
			"httpAddr": httpListenAddr,
			"grpcAddr": grpcListenAddr,
		}).Fatal("invalid env")
	}

	go func() {
		logrus.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()
	log.Fatal(makeHTTPTransport(httpListenAddr, svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	logrus.Infof("GRPC TRANSPORT RUNNING ON PORT %s", listenAddr)
	// make a tcp listener
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("error listening grpc: %v", err)
	}
	defer listener.Close()
	// make new grpc native server with options
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// register (our) grpc server implementation tothe grpc packgae
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(listener)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	aggMetricsHandler := newHTTPMetricsHandler("aggregate")
	invMetricsHandler := newHTTPMetricsHandler("invoice")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /aggregate", aggMetricsHandler.instrument(handleAggregate(svc)))
	mux.HandleFunc("GET /invoice", invMetricsHandler.instrument(handleGetInvoice(svc)))
	mux.Handle("GET /metrics", promhttp.Handler())

	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
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

		writeJSON(w, http.StatusOK, invoice)
		return
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

func makeStore(storeType string) Storer {
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		logrus.Fatalf("invalid store type given: %s", storeType)
		return nil
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}
}
