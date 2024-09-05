package main

import (
	"fmt"
	"sync"
)

// Metric types

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}

// MemStorage struct mutex map
type MemStorage struct {
	mu      sync.RWMutex
	gauge   map[string]Gauge
	counter map[string]Counter
}

// NewMemStorage ...
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]Gauge),
		counter: make(map[string]Counter),
	}
}

// Set or Update Gauge
func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauge[name] = Gauge{Name: name, Value: value}
}

// Increment Counter
func (s *MemStorage) IncCounter(name string, value int64) {
	if counter, exists := s.counter[name]; exists {
		counter.Value += value
		s.counter[name] = counter
	} else {
		s.counter[name] = Counter{Name: name, Value: value}
	}
}

// Check if POST

// Request data http://<URL>/update/<metric_type>/<metric_name>/<metric_value>
// Content-Type: text/plain

// If ok -> res http.StatusOk

// If <metric_name> == nil -> http.StatusNotFound

// ---Warning--- Not use redirects

// Just main
func main() {
	fmt.Println()
}
