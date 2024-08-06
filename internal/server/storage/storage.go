package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
)

type MemStorage struct {
	counter  sync.Map
	gauge    sync.Map
	filePath string
}

type Value struct {
	Name  string
	Type  string
	Value interface{}
}

func NewMemStorage(filePath string) *MemStorage {
	return &MemStorage{
		filePath: filePath,
	}
}

func (m *MemStorage) SetGauge(ctx context.Context, name string, value float64) error {
	m.gauge.Store(name, value)
	return nil
}

func (m *MemStorage) SetCounter(ctx context.Context, name string, value int64) error {
	if _, err := m.GetCounter(ctx, name); err != nil {
		m.counter.Store(name, value)
	} else {
		valueOld, ok := m.counter.Load(name)
		if ok {
			m.counter.Store(name, value+valueOld.(int64))
		}
	}
	return nil
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	value, ok := m.counter.Load(name)
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value.(int64), nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	value, ok := m.gauge.Load(name)
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value.(float64), nil
}

func (m *MemStorage) GetAll(ctx context.Context) ([]Value, error) {
	var values []Value

	m.gauge.Range(func(k, v interface{}) bool {
		values = append(values, Value{
			Name:  k.(string),
			Type:  "gauge",
			Value: strconv.FormatFloat(v.(float64), 'f', 1, 64),
		})
		return true
	})

	m.counter.Range(func(k, v interface{}) bool {
		values = append(values, Value{
			Name:  k.(string),
			Type:  "counter",
			Value: v.(int64),
		})
		return true
	})
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
			m.gauge.Store(metric.ID, *metric.Value)
		case "counter":
			m.counter.Store(metric.ID, *metric.Value)
		}
	}

	return nil
}

func toJSON(m *MemStorage) ([]byte, error) {
	metric := make([]metrics.Metric, 0)

	m.gauge.Range(func(k, v interface{}) bool {
		var m metrics.Metric
		m.ID = k.(string)
		m.MType = "gauge"
		newValue := v.(float64)
		m.Value = &newValue
		metric = append(metric, m)
		return true
	})

	m.counter.Range(func(k, v interface{}) bool {
		var m metrics.Metric
		m.ID = k.(string)
		m.MType = "counter"
		newDelta := v.(int64)
		m.Delta = &newDelta
		metric = append(metric, m)
		return true
	})

	return json.Marshal(metric)
}

func (m *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MemStorage) SetBatch(ctx context.Context, metrics []metrics.Metric) error {
	return nil
}
