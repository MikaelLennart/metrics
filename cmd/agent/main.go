package main

import (
	"fmt"
)

// Defalut types ...
type gauge float64
type counter int64

// Metrics ...
type Metrics struct {
	gauges      map[string]gauge
	counters    map[string]counter
	PollCount   counter
	RandomValue gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauges:      make(map[string]gauge),
		counters:    make(map[string]counter),
		PollCount:   0,
		RandomValue: 0.0,
	}
}

// Metrics Service ...
type MetricService interface {
}

// Main ...
func main() {
	fmt.Println("Boo")
}
