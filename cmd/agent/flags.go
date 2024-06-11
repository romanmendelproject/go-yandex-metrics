package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var flagReqAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&flagReqAddr, "a", "localhost:8080", "address and port to run agent")
	flag.IntVar(&reportInterval, "r", 10, "send metrics to server")
	flag.IntVar(&pollInterval, "p", 2, "collect metrics from runtime")
	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagReqAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportInterval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		}
		fmt.Println(envReportInterval)
		reportInterval = envReportInterval
	}
	fmt.Println(reportInterval)
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envPollInterval, err := strconv.Atoi(envPollInterval)
		if err != nil {
			panic(err)
		}
		pollInterval = envPollInterval
	}
}
