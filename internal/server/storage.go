package server

import (
	"fmt"
	"strconv"
	"sync"
)

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

type Storage interface {
	Update(metricType, name, value string) error
	AllGauges() map[string]float64
	AllCounters() map[string]int64
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) Update(metricType, name, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch metricType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid gauge value")
		}
		m.gauges[name] = v
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid counter value")
		}
		m.counters[name] += v
	default:
		return fmt.Errorf("invalid metric type")
	}
	return nil
}

func (m *MemStorage) AllGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		result[k] = v
	}
	return result
}

func (m *MemStorage) AllCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		result[k] = v
	}
	return result
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.gauges[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.counters[name]
	return val, ok
}
