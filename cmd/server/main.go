package main

import (
	"fmt"
	"net/http"

	"github.com/MikaelLennart/metrics.git/internal/handlers"
	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Main ...
func main() {
	storage := store.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetrics(storage))
	mux.HandleFunc("/metrics", handlers.CheckMetrics(storage))

	fmt.Println("Server started ... at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Err")
	}
}
