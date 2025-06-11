package agent

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Collector struct {
	mu         sync.RWMutex
	gauges     map[string]float64
	counters   map[string]int64
	pollCount  int64
	randomSeed *rand.Rand
}

func NewCollector() *Collector {
	return &Collector{
		gauges:     make(map[string]float64),
		counters:   make(map[string]int64),
		randomSeed: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Collector) Collect() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Сохраняем gauge-метрики
	c.gauges["Alloc"] = float64(m.Alloc)
	c.gauges["BuckHashSys"] = float64(m.BuckHashSys)
	c.gauges["Frees"] = float64(m.Frees)
	c.gauges["GCCPUFraction"] = m.GCCPUFraction
	c.gauges["GCSys"] = float64(m.GCSys)
	c.gauges["HeapAlloc"] = float64(m.HeapAlloc)
	c.gauges["HeapIdle"] = float64(m.HeapIdle)
	c.gauges["HeapInuse"] = float64(m.HeapInuse)
	c.gauges["HeapObjects"] = float64(m.HeapObjects)
	c.gauges["HeapReleased"] = float64(m.HeapReleased)
	c.gauges["HeapSys"] = float64(m.HeapSys)
	c.gauges["LastGC"] = float64(m.LastGC)
	c.gauges["Lookups"] = float64(m.Lookups)
	c.gauges["MCacheInuse"] = float64(m.MCacheInuse)
	c.gauges["MCacheSys"] = float64(m.MCacheSys)
	c.gauges["MSpanInuse"] = float64(m.MSpanInuse)
	c.gauges["MSpanSys"] = float64(m.MSpanSys)
	c.gauges["Mallocs"] = float64(m.Mallocs)
	c.gauges["NextGC"] = float64(m.NextGC)
	c.gauges["NumForcedGC"] = float64(m.NumForcedGC)
	c.gauges["NumGC"] = float64(m.NumGC)
	c.gauges["OtherSys"] = float64(m.OtherSys)
	c.gauges["PauseTotalNs"] = float64(m.PauseTotalNs)
	c.gauges["StackInuse"] = float64(m.StackInuse)
	c.gauges["StackSys"] = float64(m.StackSys)
	c.gauges["Sys"] = float64(m.Sys)
	c.gauges["TotalAlloc"] = float64(m.TotalAlloc)

	// RandomValue
	c.gauges["RandomValue"] = c.randomSeed.Float64()

	// Увеличиваем PollCount
	c.pollCount++
	c.counters["PollCount"] = c.pollCount
}

func (c *Collector) GetMetrics() (map[string]float64, map[string]int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Копируем, чтобы избежать гонок
	gauges := make(map[string]float64)
	for k, v := range c.gauges {
		gauges[k] = v
	}

	counters := make(map[string]int64)
	for k, v := range c.counters {
		counters[k] = v
	}

	return gauges, counters
}
