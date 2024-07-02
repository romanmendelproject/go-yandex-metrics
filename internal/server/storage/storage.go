package storage

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauge[name] = value
}

func (m *MemStorage) SetCounter(name string, value int64) {
	if _, err := m.GetCounter(name); err != nil {
		m.counter[name] = value
	} else {
		m.counter[name] += value
	}
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	value, ok := m.counter[name]
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value, nil
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	value, ok := m.gauge[name]
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value, nil
}

func (m *MemStorage) GetAll() []Value {
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

	return values
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

	metricSlice := make([]Metric, 0)

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
	metrics := make([]Metric, 0, len(m.gauge)+len(m.counter))

	for k, v := range m.gauge {
		var m Metric
		m.ID = k
		m.MType = "gauge"
		m.Value = &v

		metrics = append(metrics, m)
	}

	for k, v := range m.counter {
		var m Metric

		m.ID = k
		m.MType = "counter"
		m.Delta = &v
		metrics = append(metrics, m)
	}

	return json.Marshal(metrics)
}
