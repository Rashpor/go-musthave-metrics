package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

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
