package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rashpor/go-musthave-metrics/models"
	"github.com/go-chi/chi/v5"
)

func UpdateHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			http.Error(w, "Invalid request format", http.StatusNotFound)
			return
		}

		metricType, name, value := parts[0], parts[1], parts[2]

		err := storage.Update(metricType, name, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ValueHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch mType {
		case "gauge":
			val, err := storage.GetGauge(name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			fmt.Fprint(w, strconv.FormatFloat(val, 'f', -1, 64))
		case "counter":
			val, err := storage.GetCounter(name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%d", val)
		default:
			http.Error(w, "unsupported metric type", http.StatusBadRequest)
		}
	}
}

func ListHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gauges := storage.AllGauges()
		counters := storage.AllCounters()

		tmpl := `<html><body><h1>Metrics</h1><ul>
   {{ range $k, $v := .Gauges }}<li>{{$k}} (gauge): {{$v}}</li>{{ end }}
   {{ range $k, $v := .Counters }}<li>{{$k}} (counter): {{$v}}</li>{{ end }}
   </ul></body></html>`

		t := template.Must(template.New("metrics").Parse(tmpl))
		_ = t.Execute(w, map[string]interface{}{
			"Gauges":   gauges,
			"Counters": counters,
		})
	}
}

func UpdateJSONHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var m models.Metrics
		if err := json.Unmarshal(body, &m); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		switch m.MType {
		case models.Gauge:
			if m.Value == nil {
				http.Error(w, "Missing value for gauge", http.StatusBadRequest)
				return
			}
			err := storage.Update("gauge", m.ID, fmt.Sprintf("%f", *m.Value))
			if err != nil {
				http.Error(w, "Failed to update gauge", http.StatusBadRequest)
				return
			}
		case models.Counter:
			if m.Delta == nil {
				http.Error(w, "Missing delta for counter", http.StatusBadRequest)
				return
			}
			err := storage.Update("counter", m.ID, fmt.Sprintf("%d", *m.Delta))
			if err != nil {
				http.Error(w, "Failed to update counter", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown metric type", http.StatusNotImplemented)
			return
		}

		// Возвращаем обновлённую метрику
		var resp models.Metrics
		resp.ID = m.ID
		resp.MType = m.MType

		if m.MType == models.Gauge {
			val, err := storage.GetGauge(m.ID)
			if err == nil {
				resp.Value = &val
			}
		} else if m.MType == models.Counter {
			val, err := storage.GetCounter(m.ID)
			if err == nil {
				resp.Delta = &val
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func ValueJSONHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var m models.Metrics
		if err := json.Unmarshal(body, &m); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		var resp models.Metrics
		resp.ID = m.ID
		resp.MType = m.MType

		switch m.MType {
		case models.Gauge:
			val, err := storage.GetGauge(m.ID)
			if err != nil {
				http.Error(w, "Gauge not found", http.StatusNotFound)
				return
			}
			resp.Value = &val
		case models.Counter:
			val, err := storage.GetCounter(m.ID)
			if err != nil {
				http.Error(w, "Counter not found", http.StatusNotFound)
				return
			}
			resp.Delta = &val
		default:
			http.Error(w, "Unknown metric type", http.StatusNotImplemented)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
