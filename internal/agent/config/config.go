package config

import (
	"flag"
	"github.com/caarlos0/env/v8"
	"net/url"
)

type Config struct {
	JWTKey         string `env:"KEY"`
	RemoteAddr     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func GetConfiguration() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.JWTKey, "k", "", "key for jwt")
	flag.StringVar(&cfg.RemoteAddr, "a", "http://localhost:8080", "address and port of server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "flag report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "flag poll interval")
	flag.IntVar(&cfg.RateLimit, "l", 1, "rate limit")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
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
