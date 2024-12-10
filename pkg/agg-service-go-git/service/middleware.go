package aggservice

import (
	"context"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type Middleware func(Service) Service

// ### logging
type loggingMiddleware struct {
	next Service
}

func (mw *loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw *loggingMiddleware) Calculate(ctx context.Context, dist int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, dist)
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{next: next}
	}
}

// ### instrumentation
type instrumentationMiddleware struct {
	next Service
}

func (mw *instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw *instrumentationMiddleware) Calculate(ctx context.Context, dist int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, dist)
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{next: next}
	}
}
