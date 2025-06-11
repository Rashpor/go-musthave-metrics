package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func getEnvOrDefaultString(envName string, fallback string) string {
	if val := os.Getenv(envName); val != "" {
		return val
	}
	return fallback
}

func main() {

	flagAddr := flag.String("a", "", "address to run HTTP server on")
	flag.Parse()

	defaultAddr := "localhost:8080"

	addr := getEnvOrDefaultString("ADDRESS", *flagAddr)
	if addr == "" {
		addr = defaultAddr
	}

	log.Printf("Starting server on %s...", addr)
	storage := server.NewMemStorage()
	router := server.NewRouter(storage)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
