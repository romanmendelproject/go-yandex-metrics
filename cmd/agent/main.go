package main

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

		reportMetrics(&metrics)
	}
}
