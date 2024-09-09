package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/go-chi/chi/v5"
)

// Server UpdateMetrics ... Chi..v5
func UpdateMetrics(s *store.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Request: %s %s", req.Method, req.URL.Path)

		metricType := chi.URLParam(req, "type")
		metricName := chi.URLParam(req, "name")
		metricValueString := chi.URLParam(req, "value")

		switch metricType {
		case "gauge":
			log.Printf("Case gauge\r\n")
			metricValue, err := strconv.ParseFloat(metricValueString, 64)
			if err != nil {
				http.Error(res, "Value do not match", http.StatusBadRequest)
				return
			}
			s.SetGauge(metricName, metricValue)
			res.WriteHeader(http.StatusOK)

		case "counter":
			log.Printf("Case gauge\r\n")
			metricValue, err := strconv.ParseInt(metricValueString, 0, 64)
			if err != nil {
				http.Error(res, "Wrong counter value", http.StatusBadRequest)
				return
			}
			s.IncCounter(metricName, metricValue)
			res.WriteHeader(http.StatusOK)
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
		}

		if metricName == "" {
			http.Error(res, "Missing metric name", http.StatusBadRequest)
		}
	}
}
