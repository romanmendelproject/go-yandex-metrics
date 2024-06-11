package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	parseFlags()

	var metrics Metrics
	metrics.Init()

	for {
		go func() {
			err := metrics.Update()
			if err != nil {
				panic(err)
			}
		}()

		updateMetrics(&metrics)
	}
}

func updateMetrics(m *Metrics) {
	time.Sleep(time.Second * time.Duration(reportInterval))
	for k, v := range m.Data {
		if err := updateMetric(k, v); err != nil {
			panic(err)
		}
	}
}

func updateMetric(name string, metric Metric) error {
	var value interface{}

	switch metric.Value.(type) {
	case float64:
		if metric.Type != "gauge" {
			return errors.New("metric type is not float64")
		}

		valueFloat64 := metric.Value.(float64)
		value = strconv.FormatFloat(valueFloat64, 'f', 1, 64)
	case int64:
		if metric.Type != "counter" {
			return errors.New("metric type is not int64")
		}

		value = metric.Value.(int64)
	default:
		return errors.New("unknown metric type")
	}

	url := fmt.Sprintf("http://%s/update/%s/%s/%v", flagReqAddr, metric.Type, name, value)

	res, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	if err = res.Body.Close(); err != nil {
		return err
	}

	return nil
}
