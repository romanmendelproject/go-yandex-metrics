package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func reportMetrics(m *Metrics) error {
	time.Sleep(time.Second * time.Duration(reportInterval))
	for k, v := range m.DataGauge {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", flagReqAddr, v.Type, k, v.Value)
		if err := sendMetric(url); err != nil {
			log.Error(err)
		}
	}
	for k, v := range m.DataCount {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", flagReqAddr, v.Type, k, v.Value)
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
