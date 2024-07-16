package config

import (
	"flag"
	"os"
	"strconv"
)

var FlagReqAddr string
var ReportSingleInterval int
var ReportBatchInterval int
var PollInterval int

func ParseFlags() {
	flag.StringVar(&FlagReqAddr, "a", "localhost:8080", "address and port to run agent")
	flag.IntVar(&ReportSingleInterval, "r", 2, "send metrics to server")
	flag.IntVar(&ReportBatchInterval, "b", 30, "send metrics to server")
	flag.IntVar(&PollInterval, "p", 2, "collect metrics from runtime")
	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagReqAddr = envRunAddr
	}
	if envReportSingleIntervall := os.Getenv("REPORT_INTERVAL"); envReportSingleIntervall != "" {
		envReportInterval, err := strconv.Atoi(envReportSingleIntervall)
		if err != nil {
			panic(err)
		}
		ReportSingleInterval = envReportInterval
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envPollInterval, err := strconv.Atoi(envPollInterval)
		if err != nil {
			panic(err)
		}
		PollInterval = envPollInterval
	}
}
