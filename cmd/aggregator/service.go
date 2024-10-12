package main

import (
	"github.com/shariq/microservice/internal/types"
	"github.com/sirupsen/logrus"
)

type Aggregator interface {
	AggregrateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregrateDistance(distance types.Distance) error {
	logrus.Warn("processing and inserting distance in the storage: ", distance)
	return i.store.Insert(distance)
}
