// Модуль для сбора метрик с сервера
package metrics

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Metric описывает обрабатываемые метрики
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics struct {
	Data      []Metric
	PollCount int64
}

// Update получение базовых метрик
func (m *Metrics) Update(metricsChannel chan *[]Metric) error {
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)
	m.PollCount += 1

	m.Data = []Metric{
		{ID: "Alloc", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.Alloc))},
		{ID: "BuckHashSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.BuckHashSys))},
		{ID: "Frees", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.Frees))},
		{ID: "GCCPUFraction", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.GCCPUFraction))},
		{ID: "GCSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.GCSys))},
		{ID: "HeapAlloc", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapAlloc))},
		{ID: "HeapIdle", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapIdle))},
		{ID: "HeapInuse", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapInuse))},
		{ID: "HeapObjects", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapObjects))},
		{ID: "HeapReleased", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapReleased))},
		{ID: "HeapSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.HeapSys))},
		{ID: "LastGC", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.LastGC))},
		{ID: "Lookups", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.Lookups))},
		{ID: "MCacheInuse", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.MCacheInuse))},
		{ID: "MCacheSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.MCacheSys))},
		{ID: "MSpanInuse", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.MSpanInuse))},
		{ID: "MSpanSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.MSpanSys))},
		{ID: "Mallocs", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.Mallocs))},
		{ID: "NextGC", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.NextGC))},
		{ID: "NumForcedGC", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.NumForcedGC))},
		{ID: "NumGC", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.NumGC))},
		{ID: "OtherSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.OtherSys))},
		{ID: "PauseTotalNs", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.PauseTotalNs))},
		{ID: "StackInuse", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.StackInuse))},
		{ID: "StackSys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.StackSys))},
		{ID: "Sys", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.Sys))},
		{ID: "TotalAlloc", MType: "gauge", Value: utils.GetFloatPtr(float64(runtimeMetrics.TotalAlloc))},
		{ID: "RandomValue", MType: "gauge", Value: utils.GetFloatPtr(rand.Float64())},
		{ID: "PollCount", MType: "counter", Delta: &m.PollCount},
	}
	metricsChannel <- &m.Data

	return nil
}

// UpdateGopsUtil получение расширенных метрик
func (m *Metrics) UpdateGopsUtil(metricsChannel chan *[]Metric) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpuUtilMetrics, err := cpu.Percent(time.Millisecond*100, true)
	if err != nil {
		return err
	}

	var cpuUtilMetric float64
	for _, cpuUtilItem := range cpuUtilMetrics {
		cpuUtilMetric += cpuUtilItem
	}
	m.Data = []Metric{
		{ID: "TotalMemory", MType: "gauge", Value: utils.GetFloatPtr(float64(memory.Total))},
		{ID: "FreeMemory", MType: "gauge", Value: utils.GetFloatPtr(float64(memory.Free))},
		{ID: "CPUutilization1", MType: "gauge", Value: utils.GetFloatPtr(cpuUtilMetric)},
	}
	metricsChannel <- &m.Data
	return nil
}
