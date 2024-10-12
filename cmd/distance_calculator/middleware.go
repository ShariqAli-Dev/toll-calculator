package main

import (
	"time"

	"github.com/shariq/microservice/internal/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (distance float64) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"dist": distance,
		}).Info("calculated distance")
	}(time.Now())
	distance = m.next.CalculateDistance(data)
	return distance
}
