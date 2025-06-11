package main

import (
	"log"
	"net/http"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func main() {
	storage := server.NewMemStorage()
	http.HandleFunc("/update/", server.UpdateHandler(storage))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
