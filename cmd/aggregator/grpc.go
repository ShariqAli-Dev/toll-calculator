package main

import "github.com/shariqali-dev/toll-calculator/internal/types"

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) AggregateDistance(req *types.AggregateRequest) error {
	distance := types.Distance{
		OBUID:  int(req.ObuID),
		Values: req.Value,
		Unix:   req.Unix,
	}
	return s.svc.AggregateDistance(distance)
}
