package main

import (
	"flag"
)

var flagReqAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&flagReqAddr, "a", ":8080", "address and port to run agent")
	flag.IntVar(&reportInterval, "r", 10, "send metrics to server")
	flag.IntVar(&pollInterval, "p", 2, "collect metrics from runtime")
	flag.Parse()
}
