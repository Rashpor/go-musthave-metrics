package main

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := server.NewConfig()

	log.Info().Msgf("Starting server on %s...", cfg.Address)

	storage := server.NewMemStorage()
	router := server.NewRouter(storage)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}
