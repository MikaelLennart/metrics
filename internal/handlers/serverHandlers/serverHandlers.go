package serverHandlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MikaelLennart/metrics.git/internal/store"
)

func UpdateMetrics(s *store.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Request: %s %s", req.Method, req.URL.Path)
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST method", http.StatusMethodNotAllowed)
			return
		}
		urlPathString := strings.Split(req.URL.Path, "/")
		if len(urlPathString) != 5 {
			http.Error(res, "Incoreect URL", http.StatusNotFound)
			return
		}
		metricType := urlPathString[2]
		metricName := urlPathString[3]
		metricValueString := urlPathString[4]
		log.Printf("Metric Type: %s Metric Name: %s Metric Value: %s", metricType, metricName, metricValueString)
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

func CheckMetrics(s *store.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "text/plain")
		log.Println("Handler checkMetrics")
		for key, value := range s.Gauge {

			_, err := fmt.Fprintf(res, "Gauge: %s: %v\r\n", key, value)
			if err != nil {
				http.Error(res, "Error", http.StatusFailedDependency)

			}
		}

	}
}
