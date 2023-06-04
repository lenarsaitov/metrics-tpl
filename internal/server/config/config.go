package config

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

type Config struct {
	AddrRun         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

func GetConfiguration() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "time interval readings are saved to disk")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path of the file where the current values are saved")
	flag.BoolVar(&cfg.Restore, "r", true, "load previously saved values from the specified file when the server starts")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "load previously saved values from the specified file when the server starts")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
