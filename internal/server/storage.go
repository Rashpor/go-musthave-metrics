package server

import (
	"encoding/json"
	"fmt"
	"os"
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
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
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
			return fmt.Errorf("invalid gauge value: %w", err)
		}
		m.gauges[name] = v
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid counter value: %w", err)
		}
		m.counters[name] += v
	default:
		return fmt.Errorf("invalid metric type: %s", metricType)
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

func (m *MemStorage) GetGauge(name string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.gauges[name]
	if !ok {
		return 0, fmt.Errorf("gauge metric not found: %s", name)
	}
	return val, nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.counters[name]
	if !ok {
		return 0, fmt.Errorf("counter metric not found: %s", name)
	}
	return val, nil
}

// snapshot для сохранения
type snapshot struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

// SaveToFile сохраняет текущее состояние в файл
func (m *MemStorage) SaveToFile(path string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snap := snapshot{
		Gauges:   make(map[string]float64, len(m.gauges)),
		Counters: make(map[string]int64, len(m.counters)),
	}
	for k, v := range m.gauges {
		snap.Gauges[k] = v
	}
	for k, v := range m.counters {
		snap.Counters[k] = v
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(snap)
}

// LoadFromFile загружает состояние из файла, если он существует
func (m *MemStorage) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // нечего загружать
		}
		return err
	}
	defer file.Close()

	var snap snapshot
	if err := json.NewDecoder(file).Decode(&snap); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range snap.Gauges {
		m.gauges[k] = v
	}
	for k, v := range snap.Counters {
		m.counters[k] = v
	}
	return nil
}
