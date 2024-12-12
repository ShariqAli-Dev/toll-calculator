package aggservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type Middleware func(Service) Service

// ### logging
type loggingMiddleware struct {
	log  log.Logger
	next Service
}

func (mw *loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		mw.log.Log("took", time.Since(start), "obuID", dist.OBUID, "err", err)
	}(time.Now())
	err = mw.next.Aggregate(ctx, dist)
	return err
}

func (mw *loggingMiddleware) Calculate(ctx context.Context, dist int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		mw.log.Log("took", time.Since(start), "obuID", dist, "err", err)
	}(time.Now())
	invoice, err = mw.next.Calculate(ctx, dist)
	return invoice, err

}

func newLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next: next,
			log:  logger,
		}
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
