package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shariqali-dev/toll-calculator/internal/types"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg prometheus.Counter
	reqLatencyAgg prometheus.Histogram
	errCounterAgg prometheus.Counter

	reqCounterCalc prometheus.Counter
	reqLatencyCalc prometheus.Histogram
	errCounterCalc prometheus.Counter

	next Aggregator
}

type LogMiddleware struct {
	next Aggregator
}

func NewlogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	// agg
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate_counter",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate_latency",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregate_error_counter",
		Name:      "aggregate_error",
	})

	// calc
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calc_request_counter",
		Name:      "calc_counter",
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{

		Namespace: "calc_request_latency",
		Name:      "calc_latency",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calc_error_counter",
		Name:      "calc_error",
	})
	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqLatencyAgg:  reqLatencyAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyCalc: reqLatencyCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
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

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obuID)
	return invoice, err
}
