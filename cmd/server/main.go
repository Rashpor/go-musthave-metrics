package main

import (
	"log"
	"net/http"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func main() {
	cfg := server.NewConfig()

	log.Printf("Starting server on %s...", cfg.Address)

	storage := server.NewMemStorage()
	router := server.NewRouter(storage)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
