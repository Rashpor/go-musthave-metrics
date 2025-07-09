package server

import (
	"net/http"

	"github.com/Rashpor/go-musthave-metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(storage Storage) http.Handler {

	r := chi.NewRouter()

	r.Use(LoggingMiddleware)

	r.Use(middleware.DecompressMiddleware)

	r.Post("/update/{type}/{name}/{value}", UpdateHandler(storage))
	r.Post("/update/", UpdateJSONHandler(storage))
	r.Post("/value/", ValueJSONHandler(storage))
	r.Get("/value/{type}/{name}", ValueHandler(storage))
	r.Get("/", ListHandler(storage))
	return r
}
