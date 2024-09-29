package router

import (
	"net/http"

	"github.com/MikaelLennart/metrics.git/internal/handlers"
	"github.com/MikaelLennart/metrics.git/internal/middleware"
	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func NewRouter(storage *store.MemStorage, logger *logrus.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(logger))
	// POST http://<HOST>/update/<metricType>/metricName/metricName
	r.Route("/update", func(r chi.Router) {
		r.Post("/*", handlers.IsNotValidRequestURL())
		r.Post("/{type}/{name}/{value}", handlers.UpdateMetrics(storage))
	})
	// GET http://<HOST>/value/<metricType>/metricName
	r.Route("/value", func(r chi.Router) {
		r.Get("/*", handlers.IsNotValidRequestURL())
		r.Get("/{metricType}/{metricName}", handlers.GetMetricByName(storage))
	})
	// GET http://<HOST>/
	r.Get("/", handlers.GetAllMetrics(storage))

	return r
}
