package metrics

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
)

type MetricGauge struct {
	Value float64
	Type  string
}

type MetricCount struct {
	Value int64
	Type  string
}

type Metrics struct {
	DataGauge map[string]MetricGauge
	DataCount map[string]MetricCount
	PollCount int64
}

func (m *Metrics) Init() {
	m.DataGauge = make(map[string]MetricGauge)
	m.DataCount = make(map[string]MetricCount)
	m.PollCount = 0
}

func (m *Metrics) Update() error {
	time.Sleep(time.Second * time.Duration(config.PollInterval))
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.DataGauge["Alloc"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.Alloc)}
	m.DataGauge["BuckHashSys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.BuckHashSys)}
	m.DataGauge["GCCPUFraction"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.GCCPUFraction)}
	m.DataGauge["HeapAlloc"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapAlloc)}
	m.DataGauge["HeapIdle"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapIdle)}
	m.DataGauge["HeapInuse"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapInuse)}
	m.DataGauge["HeapObjects"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapObjects)}
	m.DataGauge["HeapReleased"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapReleased)}
	m.DataGauge["HeapSys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.HeapSys)}
	m.DataGauge["LastGC"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.LastGC)}
	m.DataGauge["Lookups"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.Lookups)}
	m.DataGauge["MCacheInuse"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.MCacheInuse)}
	m.DataGauge["MCacheSys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.MCacheSys)}
	m.DataGauge["MSpanInuse"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.MSpanInuse)}
	m.DataGauge["Mallocs"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.Mallocs)}
	m.DataGauge["NumForcedGC"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.NumForcedGC)}
	m.DataGauge["NumGC"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.NumGC)}
	m.DataGauge["OtherSys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.OtherSys)}
	m.DataGauge["PauseTotalNs"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.PauseTotalNs)}
	m.DataGauge["StackInuse"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.StackInuse)}
	m.DataGauge["StackSys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.StackSys)}
	m.DataGauge["Sys"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.Sys)}
	m.DataGauge["TotalAlloc"] = MetricGauge{Type: "gauge", Value: float64(runtimeMetrics.TotalAlloc)}

	m.PollCount += 1
	m.DataCount["PollCount"] = MetricCount{Type: "counter", Value: m.PollCount}
	m.DataGauge["RandomValue"] = MetricGauge{Type: "gauge", Value: rand.Float64()}

	return nil
}
