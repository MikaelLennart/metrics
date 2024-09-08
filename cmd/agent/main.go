// ~/cmd/agent
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

// Defalut types ...
type gauge float64
type counter int64

// Metrics ...
type Metrics struct {
	gauges         map[string]gauge
	counters       map[string]counter
	pollInterval   time.Duration
	reportInterval time.Duration
}

// NewMetrics ...
func NewMetrics() *Metrics {
	return &Metrics{
		gauges: map[string]gauge{
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
		counters: map[string]counter{
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
func (m *Metrics) GetMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.gauges["Alloc"] = gauge(ms.Alloc)
	m.gauges["BuckHashSys"] = gauge(ms.BuckHashSys)
	m.gauges["Frees"] = gauge(ms.Frees)
	m.gauges["GCCPUFraction"] = gauge(ms.GCCPUFraction)
	m.gauges["GCSys"] = gauge(ms.GCSys)
	m.gauges["HeapAlloc"] = gauge(ms.HeapAlloc)
	m.gauges["HeapIdle"] = gauge(ms.HeapIdle)
	m.gauges["HeapInuse"] = gauge(ms.HeapInuse)
	m.gauges["HeapObjects"] = gauge(ms.HeapObjects)
	m.gauges["HeapReleased"] = gauge(ms.HeapReleased)
	m.gauges["HeapSys"] = gauge(ms.HeapSys)
	m.gauges["LastGC"] = gauge(ms.LastGC)
	m.gauges["Lookups"] = gauge(ms.Lookups)
	m.gauges["MCacheInuse"] = gauge(ms.MCacheInuse)
	m.gauges["MCacheSys"] = gauge(ms.MCacheSys)
	m.gauges["MSpanInuse"] = gauge(ms.MSpanInuse)
	m.gauges["MSpanSys"] = gauge(ms.MSpanSys)
	m.gauges["Mallocs"] = gauge(ms.Mallocs)
	m.gauges["NextGC"] = gauge(ms.NextGC)
	m.gauges["NumForcedGC"] = gauge(ms.NumForcedGC)
	m.gauges["NumGC"] = gauge(ms.NumGC)
	m.gauges["OtherSys"] = gauge(ms.OtherSys)
	m.gauges["PauseTotalNs"] = gauge(ms.PauseTotalNs)
	m.gauges["StackInuse"] = gauge(ms.StackInuse)
	m.gauges["StackSys"] = gauge(ms.StackSys)
	m.gauges["Sys"] = gauge(ms.Sys)
	m.gauges["TotalAlloc"] = gauge(ms.TotalAlloc)
	m.gauges["RandomValue"] = gauge(rand.Float64())

	if val, exists := m.counters["PollCount"]; exists {
		m.counters["PollCount"] = val + 1
	} else {
		m.counters["PollCount"] = 1
	}
}

// MetricsPolling every 2 sec ...
func (m *Metrics) StartMetricsPolling() {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()
	for range ticker.C {
		m.GetMetrics()
	}
}

// SendMetrics by HTTP / method POST ...
func (m *Metrics) SendMetric(metricType, name string, value interface{}) {
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
func (m *Metrics) ReportMetrics() {
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()
	for range ticker.C {

		for name, value := range m.gauges {
			m.SendMetric("gauge", name, value)
			time.Sleep((100 * time.Millisecond))
		}
		for name, value := range m.counters {
			m.SendMetric("counter", name, value)
			time.Sleep((100 * time.Millisecond))
		}
	}
}

// Main ...
func main() {
	m := NewMetrics()
	go m.StartMetricsPolling()

	go m.ReportMetrics()

	select {}
}
