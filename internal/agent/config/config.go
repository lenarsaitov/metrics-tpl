package config

import (
	"flag"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	RemoteAddr     string
	ReportInterval int
	PollInterval   int
}

func GetConfiguration() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.RemoteAddr, "a", "http://localhost:8080", "address and port of server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "flag report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "flag poll interval")

	flag.Parse()

	if envRemoteAddr := os.Getenv("ADDRESS"); envRemoteAddr != "" {
		cfg.RemoteAddr = envRemoteAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		intervalSec, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return nil, err
		}
		cfg.ReportInterval = intervalSec
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		intervalSec, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return nil, err
		}
		cfg.PollInterval = intervalSec
	}

	urlRemoteAddr, err := url.Parse(cfg.RemoteAddr)
	if err != nil {
		return nil, err
	}

	if urlRemoteAddr.Scheme == "localhost" {
		cfg.RemoteAddr = "http://" + cfg.RemoteAddr
	}

	return cfg, nil
}