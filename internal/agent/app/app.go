package app

import (
	"context"
	"sync"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/report"
	"github.com/romanmendelproject/go-yandex-metrics/internal/signal"
)

func RunWorkers(ctx context.Context, wg *sync.WaitGroup, metricsChannel chan *[]metrics.Metric, workerFunc func(ctx context.Context, wg *sync.WaitGroup, metricsChannel <-chan *[]metrics.Metric)) {
	for w := 1; w <= config.RateLimit; w++ {
		wg.Add(1)
		go workerFunc(ctx, wg, metricsChannel)
	}
}

func StartAgent() {
	config.ParseFlags()

	termChan := signal.Signal()

	metricsChannel := make(chan *[]metrics.Metric, 100)
	var metr metrics.Metrics

	// tickerSingle := time.NewTicker(time.Duration(config.ReportSingleInterval) * time.Second)
	// tickerBatch := time.NewTicker(time.Duration(config.ReportBatchInterval) * time.Second)
	tickerPool := time.NewTicker(time.Duration(config.PollInterval) * time.Second)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// RunWorkers(ctx, wg, metricsChannel, report.ReportSingleMetric)
	RunWorkers(ctx, wg, metricsChannel, report.ReportBatchMetric)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerPool.C:
				go metr.Update(metricsChannel)
				go metr.UpdateGopsUtil(metricsChannel)
			}
		}
	}(ctx)

	<-termChan
	cancel()

	wg.Wait()
}
