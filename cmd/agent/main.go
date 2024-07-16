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

	tickerSingle := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	tickerBatch := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)

	for {
		go func() {
			err := metrics.Update()
			if err != nil {
				panic(err)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
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
}
