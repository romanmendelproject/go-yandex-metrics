package main

import (
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/report"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()

	var metrics metrics.Metrics

	for {
		go func() {
			err := metrics.Update()
			if err != nil {
				panic(err)
			}
		}()

		if err := report.ReportMetrics(metrics.Data); err != nil {
			log.Error(err)
		}
	}
}
