package metrics

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	metrics := Metrics{}
	metricsChannel := make(chan *[]Metric)

	go func() {
		err := metrics.Update(metricsChannel)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}()

	metricsData := <-metricsChannel

	if len(*metricsData) == 0 {
		t.Error("Expected metrics data, got none")
	}

	if metrics.PollCount != 1 {
		t.Errorf("Expected PollCount 1, got %d", metrics.PollCount)
	}

	expectedIDs := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue", "PollCount",
	}

	for _, expectedID := range expectedIDs {
		found := false
		for _, metric := range *metricsData {
			if metric.ID == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected metric ID %s not found", expectedID)
		}
	}
}

func TestUpdateGopsUtil(t *testing.T) {
	metrics := Metrics{}
	metricsChannel := make(chan *[]Metric)

	go func() {
		err := metrics.UpdateGopsUtil(metricsChannel)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}()

	metricsData := <-metricsChannel

	if len(*metricsData) == 0 {
		t.Error("Expected metrics data, got none")
	}
}
