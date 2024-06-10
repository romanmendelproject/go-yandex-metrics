package storage

import "errors"

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

func InitMemStorage() MemStorage {
	return MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauge[name] = value
}

func (m *MemStorage) SetCounter(name string, value int64) {
	if _, err := m.getCounter(name); err != nil {
		m.counter[name] = value
	} else {
		m.counter[name] += value
	}
}

func (m *MemStorage) getCounter(name string) (int64, error) {
	value, ok := m.counter[name]
	if !ok {
		return 0, errors.New("invalid name of metrics")
	}

	return value, nil
}

type Storage interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}
