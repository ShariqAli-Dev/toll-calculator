package aggendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	aggservice "github.com/shariqali-dev/toll-calculator/cmd/agg-service-go-git/service"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type AggregateRequest struct {
	types.Distance
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type AggregrateResponse struct {
	Err error `json:"err"`
}
type CalculateResponse struct {
	types.Invoice
	Err error `json:"err"`
}

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		Distance: dist,
	})
	return err
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{OBUID: obuID})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)

	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func MakeAggEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, req.Distance)
		return AggregrateResponse{
			Err: err,
		}, err
	}
}

func MakeCalcEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var invoice *types.Invoice
		req := request.(CalculateRequest)
		invoice, err = s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			Invoice: *invoice,
			Err:     err,
		}, nil
	}
}

func New(svc aggservice.Service, logger log.Logger, duration metrics.Histogram) Set {
	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = MakeAggEndpoint(svc)
		aggregateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(aggregateEndpoint)
		aggregateEndpoint = LoggingMiddleware(log.With(logger, "method", "Aggregate"))(aggregateEndpoint)
		aggregateEndpoint = InstrumentingMiddleware(duration.With("method", "Aggregate"))(aggregateEndpoint)
	}

	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = MakeAggEndpoint(svc)
		calculateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(calculateEndpoint)
		calculateEndpoint = LoggingMiddleware(log.With(logger, "method", "Calculate"))(calculateEndpoint)
		calculateEndpoint = InstrumentingMiddleware(duration.With("method", "Calculate"))(calculateEndpoint)
	}
	return Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEndpoint,
	}

}
