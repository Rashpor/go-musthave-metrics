package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func main() {
	addr := flag.String("a", "localhost:8080", "address and port to run server on")
	flag.Parse()

	storage := server.NewMemStorage()
	router := server.NewRouter(storage)

	log.Printf("Starting server on %s...\n", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
