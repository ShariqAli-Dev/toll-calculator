package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	aggendpoint "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/endpoint"
	aggservice "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/service"
	aggtransport "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/transport"
)

func main() {
	var duration metrics.Histogram
	{
		duration = kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "toll_calculator",
			Subsystem: "aggendpoint",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		service     = aggservice.New(logger)
		endpoints   = aggendpoint.New(service, logger, duration)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)

	debugListener, err := net.Listen("tcp", ":3001")
	if err != nil {
		logger.Log("transport", "debug/http", "during", "listen", "err", err)
		os.Exit(1)
	}
	go http.Serve(debugListener, http.DefaultServeMux)
	defer debugListener.Close()

	httpListener, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "listen", "err", err)
		os.Exit(1)
	}

	logger.Log("transport", "HTTP", "addr", ":3000")
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}
}
