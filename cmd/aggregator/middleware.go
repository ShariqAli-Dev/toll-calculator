package main

import (
	"time"

	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewlogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"err":   err,
			"obuID": distance.OBUID,
			"value": distance.Values,
		}).Info("AggregateDistance")
	}(time.Now())
	return m.next.AggregateDistance(distance)
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if invoice != nil {
			amount = invoice.TotalAmount
			distance = invoice.TotalDistance
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"err":         err,
			"obuID":       obuID,
			"totalDist":   distance,
			"totalAmount": amount,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obuID)
	return invoice, err
}
