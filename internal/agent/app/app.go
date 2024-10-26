// Модуль для инициализации и запуска агента
package app

import (
	"context"
	"sync"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/report"
	"github.com/romanmendelproject/go-yandex-metrics/internal/signal"

	log "github.com/sirupsen/logrus"
)

// RunWorkers запускает горутины для обработки и отпраки метрик
func RunWorkers(ctx context.Context, cfg *config.ClientFlags, wg *sync.WaitGroup, metricsChannel chan *[]metrics.Metric, workerFunc func(ctx context.Context, cfg *config.ClientFlags, wg *sync.WaitGroup, metricsChannel <-chan *[]metrics.Metric)) {
	for w := 1; w <= cfg.RateLimit; w++ {
		wg.Add(1)
		go workerFunc(ctx, cfg, wg, metricsChannel)
	}
}

// StartAgent запускает программу-агента
func StartAgent() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

	config.ParseFlags(cfg)
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

	termChan := signal.Signal()

	metricsChannel := make(chan *[]metrics.Metric, 100)
	var metr metrics.Metrics
	print(cfg.RateLimit)
	tickerPool := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunWorkers(ctx, cfg, wg, metricsChannel, report.ReportBatchMetric)

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
