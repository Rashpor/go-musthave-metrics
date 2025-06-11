package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(storage Storage) http.Handler {

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", UpdateHandler(storage))
	r.Get("/value/{type}/{name}", ValueHandler(storage))
	r.Get("/", ListHandler(storage))
	return r
}
