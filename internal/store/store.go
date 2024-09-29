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

// // Metrics Service ...
// type MetricService interface {
// 	GetMetrics()
// }

// GetMetrics ...
func (s *MemStorage) GetMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	s.Gauges["Alloc"] = Gauge(ms.Alloc)
	s.Gauges["GCCPUFraction"] = Gauge(ms.GCCPUFraction)
	s.Gauges["GCSys"] = Gauge(ms.GCSys)
	s.Gauges["Frees"] = Gauge(ms.Frees)
	s.Gauges["HeapAlloc"] = Gauge(ms.HeapAlloc)
	s.Gauges["BuckHashSys"] = Gauge(ms.BuckHashSys)
	s.Gauges["HeapIdle"] = Gauge(ms.HeapIdle)
	s.Gauges["HeapInuse"] = Gauge(ms.HeapInuse)
	s.Gauges["HeapObjects"] = Gauge(ms.HeapObjects)
	s.Gauges["HeapReleased"] = Gauge(ms.HeapReleased)
	s.Gauges["HeapSys"] = Gauge(ms.HeapSys)
	s.Gauges["LastGC"] = Gauge(ms.LastGC)
	s.Gauges["Lookups"] = Gauge(ms.Lookups)
	s.Gauges["MCacheInuse"] = Gauge(ms.MCacheInuse)
	s.Gauges["MCacheSys"] = Gauge(ms.MCacheSys)
	s.Gauges["MSpanInuse"] = Gauge(ms.MSpanInuse)
	s.Gauges["MSpanSys"] = Gauge(ms.MSpanSys)
	s.Gauges["Mallocs"] = Gauge(ms.Mallocs)
	s.Gauges["NextGC"] = Gauge(ms.NextGC)
	s.Gauges["NumForcedGC"] = Gauge(ms.NumForcedGC)
	s.Gauges["NumGC"] = Gauge(ms.NumGC)
	s.Gauges["OtherSys"] = Gauge(ms.OtherSys)
	s.Gauges["PauseTotalNs"] = Gauge(ms.PauseTotalNs)
	s.Gauges["StackInuse"] = Gauge(ms.StackInuse)
	s.Gauges["StackSys"] = Gauge(ms.StackSys)
	s.Gauges["Sys"] = Gauge(ms.Sys)
	s.Gauges["TotalAlloc"] = Gauge(ms.TotalAlloc)
	s.Gauges["RandomValue"] = Gauge(rand.Float64())

	if val, exists := s.Counters["PollCount"]; exists {
		s.Counters["PollCount"] = val + 1
	} else {
		s.Counters["PollCount"] = 1
	}
}

// MetricsPolling every <seconds> ...
func (s *MemStorage) StartMetricsPolling(seconds int64) {
	s.pollInterval = time.Duration(seconds) * time.Second
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.GetMetrics()
	}
}

// SendMetrics by HTTP / method POST ...
func (s *MemStorage) SendMetric(metricType, name string, value interface{}, address string) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", address, metricType, name, value)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte(fmt.Sprintf("%v", value))))
	if err != nil {
		fmt.Printf("Ошибка при отправке метрики: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Метрика %s %s %v отправлена. Статус: %s\n", metricType, name, value, resp.Status)
}

// ReportMetrics every <seconds> ...
func (s *MemStorage) ReportMetrics(address string, seconds int64) {
	s.reportInterval = time.Duration(seconds) * time.Second
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()
	for range ticker.C {

		for name, value := range s.Gauges {
			s.SendMetric("gauge", name, value, address)
			time.Sleep((100 * time.Millisecond))
		}
		for name, value := range s.Counters {
			s.SendMetric("counter", name, value, address)
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
