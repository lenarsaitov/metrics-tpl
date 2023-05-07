package main

import (
	"flag"
	"os"
	"strconv"
)

var flagRemoteAddr string
var flagReportInterval int
var flagPollInterval int

func parseConfiguration() error {
	flag.StringVar(&flagRemoteAddr, "a", "localhost:8080", "address and port of server")
	flag.IntVar(&flagReportInterval, "r", 10, "flag report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "flag poll interval")

	flag.Parse()

	if envRemoteAddr := os.Getenv("ADDRESS"); envRemoteAddr != "" {
		flagRemoteAddr = envRemoteAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		intervalSec, err := strconv.Atoi(envReportInterval)
		if err != nil {
			return err
		}
		flagReportInterval = intervalSec
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		intervalSec, err := strconv.Atoi(envPollInterval)
		if err != nil {
			return err
		}
		flagPollInterval = intervalSec
	}

	return nil
}
