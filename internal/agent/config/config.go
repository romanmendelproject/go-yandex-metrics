package config

import (
	"flag"
	"os"
	"strconv"
)

var FlagReqAddr string
var ReportInterval int
var PollInterval int

func ParseFlags() {
	flag.StringVar(&FlagReqAddr, "a", "localhost:8080", "address and port to run agent")
	flag.IntVar(&ReportInterval, "r", 10, "send metrics to server")
	flag.IntVar(&PollInterval, "p", 2, "collect metrics from runtime")
	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagReqAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportInterval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
		}
		ReportInterval = envReportInterval
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envPollInterval, err := strconv.Atoi(envPollInterval)
		if err != nil {
			panic(err)
		}
		PollInterval = envPollInterval
	}
}
