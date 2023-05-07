package main

import (
	"flag"
	"os"
)

var flagAddrRun string

func parseConfiguration() {
	flag.StringVar(&flagAddrRun, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envAddrRun := os.Getenv("ADDRESS"); envAddrRun != "" {
		flagAddrRun = envAddrRun
	}
}
