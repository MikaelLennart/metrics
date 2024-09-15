package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

// Config struct ...
type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	PollInterval         int64  `env:"POLL_INTERVAL" envDefault:"2"`
	ReportInterval       int64  `env:"REPORT_INTERVAL" envDefault:"10"`
	ServerAddressEnvSet  bool
	PollIntervalEnvSet   bool
	ReportIntervalEnvSet bool
}

// checkEnvSet ... check if env is set
func checkEnvSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// ServerConfig ...
func ServerConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	cfg.ServerAddressEnvSet = checkEnvSet("SERVER_ADDRESS")
	cfg.PollIntervalEnvSet = checkEnvSet("POLL_INTERVAL")
	cfg.ReportIntervalEnvSet = checkEnvSet("REPORT_INTERVAL")

	if !cfg.ServerAddressEnvSet {
		serverAddress := flag.String("a", cfg.ServerAddress, "Server address")
		flag.Parse()
		cfg.ServerAddress = *serverAddress
		fmt.Println(*serverAddress)
	}

	return cfg

}

func AgentConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if !cfg.ServerAddressEnvSet {
		serverAddress := flag.String("a", cfg.ServerAddress, "Server address")
		flag.Parse()
		cfg.ServerAddress = *serverAddress
		fmt.Println(*serverAddress)
	}

	if !cfg.PollIntervalEnvSet {
		pollInterval := flag.Int64("p", cfg.PollInterval, "Set Poll Interval in Seconds")
		flag.Parse()
		cfg.PollInterval = *pollInterval
		fmt.Println(*pollInterval)
	}

	if !cfg.ReportIntervalEnvSet {
		reportInterval := flag.Int64("p", cfg.ReportInterval, "Set Report Interval in Seconds")
		flag.Parse()
		cfg.ReportInterval = *reportInterval
		fmt.Println(*reportInterval)
	}

	return cfg

}
