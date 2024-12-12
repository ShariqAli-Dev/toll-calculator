package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	aggendpoint "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/endpoint"
	aggservice "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/service"
	aggtransport "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/transport"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var (
		service     = aggservice.New()
		endpoints   = aggendpoint.New(service, logger)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)

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
