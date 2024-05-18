package main

import (
	"time"

	types "github.com/shariqali-dev/toll-calculator/internal"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewlogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info("calculate distance")
	}(time.Now())

	dist, err = m.next.CalculateDistance(data)
	return dist, err
}
