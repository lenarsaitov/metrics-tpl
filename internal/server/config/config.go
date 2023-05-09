package config

import (
	"flag"
	"net/url"
	"os"
)

type Config struct {
	AddrRun string
}

func GetConfiguration() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envAddrRun := os.Getenv("ADDRESS"); envAddrRun != "" {
		cfg.AddrRun = envAddrRun
	}

	_, err := url.Parse(cfg.AddrRun)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
