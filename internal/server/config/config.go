package config

import (
	"flag"
	"os"
)

type Config struct {
	AddrRun string
}

func GetConfiguration() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envAddrRun := os.Getenv("ADDRESS"); envAddrRun != "" {
		cfg.AddrRun = envAddrRun
	}

	return cfg
}
