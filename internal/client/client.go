package client

import (
	"context"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
