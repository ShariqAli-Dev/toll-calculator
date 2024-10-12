package main

import (
	"math"

	"github.com/shariq/microservice/internal/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) float64
}

type CalculatorService struct {
	prevPoint *types.OBUData
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) float64 {
	if s.prevPoint == nil {
		s.prevPoint = &data
		return 0
	}

	distance := calculateDistance(data.Lat, s.prevPoint.Lat, data.Long, s.prevPoint.Long)
	s.prevPoint = &data
	return distance
}

func calculateDistance(x2, x1, y2, y1 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
