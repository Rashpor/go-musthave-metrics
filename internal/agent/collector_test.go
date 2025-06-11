package agent_test

import (
	"testing"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	c := agent.NewCollector()

	// Собираем метрики
	c.Collect()

	// Получаем собранные метрики
	gauges, counters := c.GetMetrics()

	require.Greater(t, len(gauges), 0, "gauges should not be empty")
	require.Greater(t, len(counters), 0, "counters should not be empty")

	// Проверяем наличие конкретной метрики
	_, ok := gauges["HeapAlloc"]
	require.True(t, ok, "HeapAlloc metric should be present")

	// Проверяем, что PollCount увеличился
	val, ok := counters["PollCount"]
	require.True(t, ok, "PollCount should be present")
	require.Equal(t, int64(1), val, "PollCount should be 1 after first collection")
}
