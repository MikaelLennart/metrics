package main

import (
	"fmt"
	"net/http"

	"github.com/MikaelLennart/metrics.git/internal/handlers"
	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/go-chi/chi/v5"
)

// Main ...
func main() {
	storage := store.NewMemStorage()
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics(storage))

	fmt.Println("Server started ... at :8080")
	http.ListenAndServe(":8080", r)
}
