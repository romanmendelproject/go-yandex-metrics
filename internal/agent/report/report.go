package report

import (
	"fmt"
	"net/http"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"

	log "github.com/sirupsen/logrus"
)

func ReportMetrics(m *metrics.Metrics) error {
	time.Sleep(time.Second * time.Duration(config.ReportInterval))
	for k, v := range m.DataGauge {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", config.FlagReqAddr, v.Type, k, v.Value)
		if err := sendMetric(url); err != nil {
			log.Error(err)
		}
	}
	for k, v := range m.DataCount {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", config.FlagReqAddr, v.Type, k, v.Value)
		if err := sendMetric(url); err != nil {
			log.Error(err)
		}
	}

	return nil
}

func sendMetric(url string) error {
	res, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	if err = res.Body.Close(); err != nil {
		return err
	}

	return nil
}
