package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/shariqali-dev/toll-calculator/internal/client"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3001", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3002", "the listen address of the GRPC server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewlogMiddleware(NewInvoiceAggregator(store))

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc))
	}()
	time.Sleep(time.Second * 2)
	grpcClient, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = grpcClient.Client.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 50.8,
		Unix:  time.Now().Unix(),
	}); err != nil {
		log.Fatal(err)
	}
	log.Fatal(makeHTTPTransport(*httpListenAddr, svc))
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
	mux := http.NewServeMux()
	mux.Handle("POST /aggregate", handleAggregate(svc))
	mux.Handle("GET /invoice", handleGetInvoice(svc))

	logrus.Infof("HTTP TRANSPORT RUNNING ON PORT %s", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithField("thes tring", r.URL.Query().Get("obuID")).Info("THE OBU ID ERRORING")
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
