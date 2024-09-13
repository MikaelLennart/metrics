// ~/cmd/agent
package main

import (
	"flag"

	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Agent Main ...
func main() {
	address := flag.String("a", "localhost:8080", "server port adress")
	reportInterval := flag.Int64("r", 10, "metrics report to server interval")
	pollInterval := flag.Int64("p", 2, "metrics polling interval")
	flag.Parse()
	port := "" + *address
	s := store.NewMetrics()

	go s.StartMetricsPolling(*pollInterval)
	go s.ReportMetrics(port, *reportInterval)

	select {}
}
