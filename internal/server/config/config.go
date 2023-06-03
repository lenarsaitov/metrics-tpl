package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	AddrRun         string
	FileStoragePath string
	DatabaseDSN     string
	StoreInterval   int
	Restore         bool
}

func GetConfiguration() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.AddrRun, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "time interval readings are saved to disk")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path of the file where the current values are saved")
	flag.BoolVar(&cfg.Restore, "r", true, "load previously saved values from the specified file when the server starts")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "load previously saved values from the specified file when the server starts")

	flag.Parse()

	if envAddrRun := os.Getenv("ADDRESS"); envAddrRun != "" {
		cfg.AddrRun = envAddrRun
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		number, err := strconv.Atoi(envStoreInterval)
		if nil == err {
			cfg.StoreInterval = number
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		boolean, err := strconv.ParseBool(envRestore)
		if nil == err {
			cfg.Restore = boolean
		}
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	return cfg
}
