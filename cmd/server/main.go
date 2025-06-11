package main

import (
	"log"
	"net/http"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func main() {
	storage := server.NewMemStorage()
	router := server.NewRouter(storage)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
