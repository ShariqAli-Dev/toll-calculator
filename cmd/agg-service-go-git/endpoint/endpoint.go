package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	aggservice "github.com/shariqali-dev/toll-calculator/pkg/agg-service-go-git/service"
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
