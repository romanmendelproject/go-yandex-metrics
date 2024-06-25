package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"

	log "github.com/sirupsen/logrus"
)

func ReportMetrics(data []metrics.Metric) error {
	time.Sleep(time.Second * time.Duration(config.ReportInterval))
	for _, v := range data {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			log.Error(err)
		} else {
			if err := sendMetric(jsonValue); err != nil {
				log.Error(err)
			}
		}
	}

	return nil
}

func sendMetric(body []byte) error {
	url := fmt.Sprintf("http://%s/update/", config.FlagReqAddr)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if err = res.Body.Close(); err != nil {
		return err
	}

	return nil
}
