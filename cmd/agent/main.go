package main

import log "github.com/sirupsen/logrus"

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

		if err := reportMetrics(&metrics); err != nil {
			log.Error(err)
		}
	}
}
