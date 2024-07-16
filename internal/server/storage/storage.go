package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
)

type MemStorage struct {
	counter  map[string]int64
	gauge    map[string]float64
	filePath string
}

type Value struct {
	Name  string
	Type  string
	Value interface{}
}

func NewMemStorage(filePath string) *MemStorage {
	return &MemStorage{
		gauge:    make(map[string]float64),
		counter:  make(map[string]int64),
		filePath: filePath,
	}
}

func (m *MemStorage) SetGauge(ctx context.Context, name string, value float64) error {
	m.gauge[name] = value
	return nil
}

func (m *MemStorage) SetCounter(ctx context.Context, name string, value int64) error {
	if _, err := m.GetCounter(ctx, name); err != nil {
		m.counter[name] = value
	} else {
		m.counter[name] += value
	}
	return nil
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	value, ok := m.counter[name]
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value, nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	value, ok := m.gauge[name]
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value, nil
}

func (m *MemStorage) GetAll(ctx context.Context) ([]Value, error) {
	var values []Value

	for k, v := range m.gauge {
		values = append(values, Value{
			Name:  k,
			Type:  "gauge",
			Value: strconv.FormatFloat(v, 'f', 1, 64),
		})
	}

	for k, v := range m.counter {
		values = append(values, Value{
			Name:  k,
			Type:  "counter",
			Value: v,
		})
	}
	fmt.Println(values)
	return values, nil
}

func (m *MemStorage) SaveToFile() error {
	file, err := os.Create(m.filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	metrics, err := toJSON(m)
	if err != nil {
		return err
	}

	_, err = file.Write(metrics)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemStorage) RestoreFromFile() error {
	file, err := os.ReadFile(m.filePath)

	if err != nil {
		return err
	}

	metricSlice := make([]metrics.Metric, 0)

	err = json.Unmarshal(file, &metricSlice)
	if err != nil {
		return err
	}

	for _, metric := range metricSlice {
		switch metric.MType {
		case "gauge":
			m.gauge[metric.ID] = *metric.Value
		case "counter":
			m.counter[metric.ID] = *metric.Delta
		}
	}

	return nil
}

func toJSON(m *MemStorage) ([]byte, error) {
	metric := make([]metrics.Metric, 0, len(m.gauge)+len(m.counter))

	for k, v := range m.gauge {
		var m metrics.Metric
		m.ID = k
		m.MType = "gauge"
		newValue := v
		m.Value = &newValue

		metric = append(metric, m)
	}

	for k, v := range m.counter {
		var m metrics.Metric
		m.ID = k
		m.MType = "counter"
		newDelta := v
		m.Delta = &newDelta
		metric = append(metric, m)
	}

	return json.Marshal(metric)
}

func (m *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemStorage) SetBatch(ctx context.Context, metrics []metrics.Metric) error {
	return nil
}
