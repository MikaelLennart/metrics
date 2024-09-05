package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

// Update Metrtics
// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (s *Server) UpdateMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST method", http.StatusMethodNotAllowed)
	}
	urlPathString := strings.Split(req.URL.Path, "/")
	if len(urlPathString) != 5 {
		http.Error(res, "Incoreect URL", http.StatusNotFound)
	} else {
		io.WriteString(res, "All good\r\n")
	}

	metricType := urlPathString[2]

	if metricType != "gauge" && metricType != "counter" {
		http.Error(res, "Wrong metric type", http.StatusNotFound)
	} else {
		res.WriteHeader(http.StatusOK)
	}
	// metricName := urlPathString[4]
	// metricValue := urlPathString[5]

}

// Request data http://<URL>/update/<metric_type>/<metric_name>/<metric_value>
// Content-Type: text/plain

// If ok -> res http.StatusOk

// If <metric_name> == nil -> http.StatusNotFound

// ---Warning--- Not use redirects

type Server struct {
	storage *MemStorage
}

func NewServer(storage *MemStorage) *Server {
	return &Server{storage: storage}
}

// Just main
func main() {
	storage := NewMemStorage()
	server := NewServer(storage)
	fmt.Println()
	http.HandleFunc("/update/", server.UpdateMetrics)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Err")
	}
}
