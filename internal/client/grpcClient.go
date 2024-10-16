package client

import (
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint:         endpoint,
		AggregatorClient: client,
	}, nil
}
