package main

import (
	"context"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID:  int(req.ObuID),
		Values: req.Value,
		Unix:   req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
