package main

import (
	"math/rand"
	"runtime"
	"time"
)

type Metric struct {
	Type  string
	Value interface{}
}

type Metrics struct {
	Data      map[string]Metric
	PollCount int64
}

func (m *Metrics) Init() {
	m.Data = make(map[string]Metric)
	m.PollCount = 0
}

func (m *Metrics) Update() error {
	time.Sleep(time.Second * time.Duration(pollInterval))
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.Data["Alloc"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.Alloc)}
	m.Data["BuckHashSys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.BuckHashSys)}
	m.Data["GCCPUFraction"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.GCCPUFraction)}
	m.Data["HeapAlloc"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapAlloc)}
	m.Data["HeapIdle"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapIdle)}
	m.Data["HeapInuse"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapInuse)}
	m.Data["HeapObjects"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapObjects)}
	m.Data["HeapReleased"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapReleased)}
	m.Data["HeapSys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.HeapSys)}
	m.Data["LastGC"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.LastGC)}
	m.Data["Lookups"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.Lookups)}
	m.Data["MCacheInuse"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.MCacheInuse)}
	m.Data["MCacheSys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.MCacheSys)}
	m.Data["MSpanInuse"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.MSpanInuse)}
	m.Data["Mallocs"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.Mallocs)}
	m.Data["NumForcedGC"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.NumForcedGC)}
	m.Data["NumGC"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.NumGC)}
	m.Data["OtherSys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.OtherSys)}
	m.Data["PauseTotalNs"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.PauseTotalNs)}
	m.Data["StackInuse"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.StackInuse)}
	m.Data["StackSys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.StackSys)}
	m.Data["Sys"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.Sys)}
	m.Data["TotalAlloc"] = Metric{Type: "gauge", Value: float64(runtimeMetrics.TotalAlloc)}

	m.PollCount += 1
	m.Data["PollCount"] = Metric{Type: "counter", Value: m.PollCount}
	m.Data["RandomValue"] = Metric{Type: "gauge", Value: rand.Float64()}

	return nil
}
