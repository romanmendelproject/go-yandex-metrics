package main

import (
	"context"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/report"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()

	var metrics metrics.Metrics

	ctx := context.Background()

	tickerSingle := time.NewTicker(time.Duration(config.ReportSingleInterval) * time.Second)
	tickerBatch := time.NewTicker(time.Duration(config.ReportBatchInterval) * time.Second)
	tickerPool := time.NewTicker(time.Duration(config.PollInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tickerPool.C:
			err := metrics.Update()
			if err != nil {
				panic(err)
			}
		case <-tickerSingle.C:
			if err := report.ReportSingleMetric(metrics.Data); err != nil {
				log.Error(err)
			}

		case <-tickerBatch.C:
			if err := report.ReportBatchMetrics(metrics.Data); err != nil {
				log.Error(err)
			}
		}
	}
}
