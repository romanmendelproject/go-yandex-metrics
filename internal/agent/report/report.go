// Модуль отправки данных на сервер
package report

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/crypto"

	log "github.com/sirupsen/logrus"
)

var retries = []int{1, 3, 5}

// ReportSingleMetric отправка одинарной метрики на сервер
func ReportSingleMetric(ctx context.Context, wg *sync.WaitGroup, metricsChannel <-chan *[]metrics.Metric) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Debug("Stoping single sender gorutine")
			return
		case data := <-metricsChannel:
			for _, v := range *data {
				jsonValue, err := json.Marshal(v)
				if err != nil {
					log.Error(err)
				} else {
					url := fmt.Sprintf("http://%s/update/", config.FlagReqAddr)
					if err := sendMetric(jsonValue, url); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

// ReportBatchMetric отправка нескольких метрик в одном пакете в формате JSON
func ReportBatchMetric(ctx context.Context, wg *sync.WaitGroup, metricsChannel <-chan *[]metrics.Metric) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Debug("Stoping batch sender gorutine")
			return
		case data := <-metricsChannel:
			jsonValue, err := json.Marshal(data)
			if err != nil {
				log.Error(err)
			} else {
				url := fmt.Sprintf("http://%s/updates/", config.FlagReqAddr)
				if err := sendMetric(jsonValue, url); err != nil {
					log.Error(err)
				}
			}

		}
	}
}

func sendMetric(body []byte, url string) error {
	var requestBody bytes.Buffer

	gz := gzip.NewWriter(&requestBody)
	gz.Write(body)
	gz.Close()

	client := http.Client{}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Error(err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if config.Key != "" {
		hash := crypto.GetHash(body, config.Key)
		req.Header.Set("HashSHA256", hash)
	}

	for _, timeSleep := range retries {
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Failed to send collectors to server: %s. Retrying after %ds...", err, timeSleep)
			time.Sleep(time.Duration(timeSleep) * time.Second)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("not expected status code: %d", resp.StatusCode)
		} else {
			return nil
		}
	}

	return err
}
