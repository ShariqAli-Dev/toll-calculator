package main

import types "github.com/shariqali-dev/toll-calculator/internal"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	return nil
}
