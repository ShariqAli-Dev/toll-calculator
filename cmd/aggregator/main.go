package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shariqali-dev/toll-calculator/internal/store"
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

	aggHandler := makeHTTPHandlerFunc(aggMetricsHandler.instrument(handleAggregate(svc)))
	invHandler := makeHTTPHandlerFunc(invMetricsHandler.instrument(handleGetInvoice(svc)))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /invoice", invHandler)
	mux.HandleFunc("POST /aggregate", aggHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

func handleGetInvoice(service Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		obuID, err := strconv.Atoi(r.URL.Query().Get("obuID"))
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing or invalid obuID"),
			}
		}

		invoice, err := service.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, distance)

	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeStore(storeType string) store.Storer {
	switch storeType {
	case "memory":
		return store.NewMemoryStore()
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
