package report

import (
	"bytes"
	"compress/gzip"
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

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return nil
}
