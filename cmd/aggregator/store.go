package main

import (
	"fmt"

	"github.com/shariqali-dev/toll-calculator/internal/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[int]float64{},
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Values
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	distance, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id %d", id)
	}
	return distance, nil
}
