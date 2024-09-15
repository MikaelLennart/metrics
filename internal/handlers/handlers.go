package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/go-chi/chi/v5"
)

// Server UpdateMetrics ... Chi..v5
func UpdateMetrics(s *store.MemStorage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// log.Printf("Request: %s %s", req.Method, req.URL.Path)

		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValueString := chi.URLParam(r, "value")

		switch metricType {
		case "gauge":
			// log.Printf("#gauge\r\n")
			metricValue, err := strconv.ParseFloat(metricValueString, 64)
			if err != nil {
				http.Error(rw, "Value do not match", http.StatusBadRequest)
				return
			}
			s.SetGauge(metricName, metricValue)
			rw.WriteHeader(http.StatusOK)
			// log.Printf("200 : StatusOK")
		case "counter":
			// log.Printf("#counter\r\n")
			metricValue, err := strconv.ParseInt(metricValueString, 0, 64)
			if err != nil {
				http.Error(rw, "Wrong counter value", http.StatusBadRequest)
				return
			}
			s.IncCounter(metricName, metricValue)
			rw.WriteHeader(http.StatusOK)
			// log.Printf("200 : StatusOK")
		default:
			http.Error(rw, "Wrong metric type", http.StatusBadRequest)
		}

	}
}

// Get Metric by name...
func GetMetricByName(s *store.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		switch metricType {
		case "gauge":
			if _, exists := s.Gauges[metricName]; exists {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%v\r\n", s.Gauges[metricName])))
			} else {
				http.Error(w, "Metric not fount", http.StatusNotFound)
			}
		case "counter":
			if _, exists := s.Counters[metricName]; exists {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%v\r\n", s.Counters[metricName])))
			} else {
				http.Error(w, "Metric not fount", http.StatusNotFound)
			}

		}
	}
}

// GetMetrics ...
func GetAllMetrics(s *store.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "Gauge metrics:\r\n")
		for key, value := range s.Gauges {
			w.Write([]byte(fmt.Sprintf("Key: %v Value: %v\r\n", key, value)))
		}
		io.WriteString(w, "Countetr metrics:\r\n")
		for key, value := range s.Counters {
			w.Write([]byte(fmt.Sprintf("Key: %v Value: %v\r\n", key, value)))
		}

	}
}

// Router validation ...
func IsNotValidRequestURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
	}
}
