package storage

import (
	"errors"
	"strconv"
)

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

type Value struct {
	Name  string
	Type  string
	Value interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
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
