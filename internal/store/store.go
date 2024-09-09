package store

import (
	"sync"
)

// MemStorage struct mutex map
type MemStorage struct {
	mu      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

// NewMemStorage ...
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// Set or Update Gauge
func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gauge[name] = value
}

// Increment Counter
func (s *MemStorage) IncCounter(name string, value int64) {
	if counter, exists := s.Counter[name]; exists {
		s.Counter[name] = counter + value
	} else {
		s.Counter[name] = value
	}
}
