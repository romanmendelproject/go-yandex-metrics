package main

import (
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()

	var metrics metrics.Metrics
	metrics.Init()

	for {
		go func() {
			err := metrics.Update()
			if err != nil {
				panic(err)
			}
		}()

		if err := report.ReportMetrics(&metrics); err != nil {
			log.Error(err)
		}
	}
}
