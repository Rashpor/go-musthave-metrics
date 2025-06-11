package agent_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestSender_sendMetric(t *testing.T) {
	var received bool

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/update/gauge/test/1.23" {
			received = true
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer ts.Close()

	sender := agent.NewSender(ts.URL)

	gauges := map[string]float64{
		"test": 1.23,
	}
	counters := map[string]int64{}

	err := sender.Send(gauges, counters)
	require.NoError(t, err)
	require.True(t, received, "Metric was not received by test server")
}
