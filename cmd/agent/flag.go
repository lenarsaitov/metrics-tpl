package main

import "flag"

var flagRemoteAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlag() {
	flag.StringVar(&flagRemoteAddr, "a", "localhost:8080", "address and port of server")
	flag.IntVar(&flagReportInterval, "r", 10, "flag report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "flag poll interval")

	flag.Parse()
}
