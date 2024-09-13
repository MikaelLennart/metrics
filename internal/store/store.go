package store

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Defalut types ...
type Gauge float64
type Counter int64

// Metrics ...
type MemStorage struct {
	mu             sync.RWMutex
	Gauges         map[string]Gauge
	Counters       map[string]Counter
	pollInterval   time.Duration
	reportInterval time.Duration
}

// NewMemStorage ...
func NewMemStorage() *MemStorage {
	return &MemStorage{

		Gauges:   make(map[string]Gauge),
		Counters: make(map[string]Counter),
	}
}

// NewMetrics ...
func NewMetrics() *MemStorage {
	return &MemStorage{
		Gauges: map[string]Gauge{
			"Alloc":         0,
			"BuckHashSys":   0,
			"Frees":         0,
			"GCCPUFraction": 0,
			"GCSys":         0,
			"HeapAlloc":     0,
			"HeapIdle":      0,
			"HeapInuse":     0,
			"HeapObjects":   0,
			"HeapReleased":  0,
			"HeapSys":       0,
			"LastGC":        0,
			"Lookups":       0,
			"MCacheInuse":   0,
			"MCacheSys":     0,
			"MSpanInuse":    0,
			"MSpanSys":      0,
			"Mallocs":       0,
			"NextGC":        0,
			"NumForcedGC":   0,
			"NumGC":         0,
			"OtherSys":      0,
			"PauseTotalNs":  0,
			"StackInuse":    0,
			"StackSys":      0,
			"Sys":           0,
			"TotalAlloc":    0,
			"RandomValue":   0.0,
		},
		Counters: map[string]Counter{
			"PollCount": 0,
		},
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
	}
}

// Metrics Service ...
type MetricService interface {
	GetMetrics()
}

// GetMetrics ...
func (m *MemStorage) GetMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.Gauges["Alloc"] = Gauge(ms.Alloc)
	m.Gauges["GCCPUFraction"] = Gauge(ms.GCCPUFraction)
	m.Gauges["GCSys"] = Gauge(ms.GCSys)
	m.Gauges["Frees"] = Gauge(ms.Frees)
	m.Gauges["HeapAlloc"] = Gauge(ms.HeapAlloc)
	m.Gauges["BuckHashSys"] = Gauge(ms.BuckHashSys)
	m.Gauges["HeapIdle"] = Gauge(ms.HeapIdle)
	m.Gauges["HeapInuse"] = Gauge(ms.HeapInuse)
	m.Gauges["HeapObjects"] = Gauge(ms.HeapObjects)
	m.Gauges["HeapReleased"] = Gauge(ms.HeapReleased)
	m.Gauges["HeapSys"] = Gauge(ms.HeapSys)
	m.Gauges["LastGC"] = Gauge(ms.LastGC)
	m.Gauges["Lookups"] = Gauge(ms.Lookups)
	m.Gauges["MCacheInuse"] = Gauge(ms.MCacheInuse)
	m.Gauges["MCacheSys"] = Gauge(ms.MCacheSys)
	m.Gauges["MSpanInuse"] = Gauge(ms.MSpanInuse)
	m.Gauges["MSpanSys"] = Gauge(ms.MSpanSys)
	m.Gauges["Mallocs"] = Gauge(ms.Mallocs)
	m.Gauges["NextGC"] = Gauge(ms.NextGC)
	m.Gauges["NumForcedGC"] = Gauge(ms.NumForcedGC)
	m.Gauges["NumGC"] = Gauge(ms.NumGC)
	m.Gauges["OtherSys"] = Gauge(ms.OtherSys)
	m.Gauges["PauseTotalNs"] = Gauge(ms.PauseTotalNs)
	m.Gauges["StackInuse"] = Gauge(ms.StackInuse)
	m.Gauges["StackSys"] = Gauge(ms.StackSys)
	m.Gauges["Sys"] = Gauge(ms.Sys)
	m.Gauges["TotalAlloc"] = Gauge(ms.TotalAlloc)
	m.Gauges["RandomValue"] = Gauge(rand.Float64())

	if val, exists := m.Counters["PollCount"]; exists {
		m.Counters["PollCount"] = val + 1
	} else {
		m.Counters["PollCount"] = 1
	}
}

// MetricsPolling every 2 sec ...
func (s *MemStorage) StartMetricsPolling() {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.GetMetrics()
	}
}

// SendMetrics by HTTP / method POST ...
func (s *MemStorage) SendMetric(metricType, name string, value interface{}) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, name, value)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte(fmt.Sprintf("%v", value))))
	if err != nil {
		fmt.Printf("Ошибка при отправке метрики: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Метрика %s %s %v отправлена. Статус: %s\n", metricType, name, value, resp.Status)
}

// ReportMetrics every 5 sec ..
func (s *MemStorage) ReportMetrics() {
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()
	for range ticker.C {

		for name, value := range s.Gauges {
			s.SendMetric("gauge", name, value)
			time.Sleep((100 * time.Millisecond))
		}
		for name, value := range s.Counters {
			s.SendMetric("counter", name, value)
			time.Sleep((100 * time.Millisecond))
		}
	}
}

// Set or Update Gauge
func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gauges[name] = Gauge(value)
}

// Increment Counter
func (s *MemStorage) IncCounter(name string, value int64) {
	if counter, exists := s.Counters[name]; exists {
		s.Counters[name] = Counter(counter) + Counter(value)
	} else {
		s.Counters[name] = Counter(value)
	}
}

//
