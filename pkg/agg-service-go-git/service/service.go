package aggservice

import (
	"context"

	"github.com/shariqali-dev/toll-calculator/internal/store"
	"github.com/shariqali-dev/toll-calculator/internal/types"
)

const basePrice = 3.15

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store store.Storer
}

func newBasicService(store store.Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return inv, nil
}

func New() Service {
	var svc Service
	{
		svc = newBasicService(store.NewMemoryStore())
		svc = newLoggingMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}
	return svc
}
