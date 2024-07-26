package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/report"
)

func RunSenderSingleWorkers(ctx context.Context, wg *sync.WaitGroup, metricsChannel chan *[]metrics.Metric) {
	for w := 1; w <= config.RateLimit; w++ {
		wg.Add(1)
		go report.ReportSingleMetric(ctx, wg, metricsChannel)
	}
}

func RunSenderBatchWorkers(ctx context.Context, wg *sync.WaitGroup, metricsChannel chan *[]metrics.Metric) {
	for w := 1; w <= config.RateLimit; w++ {
		wg.Add(1)
		go report.ReportBatchMetric(ctx, wg, metricsChannel)
	}
}

func main() {
	config.ParseFlags()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	metricsChannel := make(chan *[]metrics.Metric, 100)
	var metr metrics.Metrics

	// tickerSingle := time.NewTicker(time.Duration(config.ReportSingleInterval) * time.Second)
	// tickerBatch := time.NewTicker(time.Duration(config.ReportBatchInterval) * time.Second)
	tickerPool := time.NewTicker(time.Duration(config.PollInterval) * time.Second)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunSenderSingleWorkers(ctx, wg, metricsChannel)
	RunSenderBatchWorkers(ctx, wg, metricsChannel)

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
