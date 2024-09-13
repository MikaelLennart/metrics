// ~/cmd/agent
package main

import (
	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Agent Main ...
func main() {
	s := store.NewMetrics()

	go s.StartMetricsPolling()
	go s.ReportMetrics()

	select {}
}
