// ~/cmd/agent
package main

import (
	"github.com/MikaelLennart/metrics.git/config"
	"github.com/MikaelLennart/metrics.git/internal/store"
)

// Agent Main ...
func main() {
	cfg := config.AgentConfig()
	s := store.NewMetrics()

	go s.StartMetricsPolling(cfg.PollInterval)
	go s.ReportMetrics(cfg.ServerAddress, cfg.ReportInterval)

	select {}
}
