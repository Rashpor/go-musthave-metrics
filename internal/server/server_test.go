package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func TestUpdateHandler(t *testing.T) {
	storage := server.NewMemStorage()
	handler := server.UpdateHandler(storage)

	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{"valid gauge", http.MethodPost, "/update/gauge/testGauge/123.45", http.StatusOK},
		{"valid counter", http.MethodPost, "/update/counter/testCounter/5", http.StatusOK},
		{"missing metric name", http.MethodPost, "/update/counter//5", http.StatusNotFound},
		{"invalid metric type", http.MethodPost, "/update/invalid/test/1", http.StatusBadRequest},
		{"invalid value", http.MethodPost, "/update/counter/test/invalid", http.StatusBadRequest},
		{"wrong method", http.MethodGet, "/update/counter/test/1", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Set("Content-Type", "text/plain")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
