package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Rashpor/go-musthave-metrics/internal/server"
)

/*
func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := server.NewConfig()

	log.Info().Msgf("Starting server on %s...", cfg.Address)

	storage := server.NewMemStorage()
	router := server.NewRouter(storage)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}*/

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := server.NewConfig()
	log.Info().Msgf("Starting server on %s...", cfg.Address)
	log.Info().Msgf("StoreInterval=%s, File=%q, Restore=%v", cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)

	storage := server.NewMemStorage()

	// 1) Восстановление
	if cfg.Restore {
		if err := storage.LoadFromFile(cfg.FileStoragePath); err != nil {
			log.Error().Err(err).Msg("failed to restore metrics from file")
		} else {
			log.Info().Msg("restored metrics from file")
		}
	}

	router := server.NewRouter(storage)
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// 2) Периодическое сохранение
	done := make(chan struct{})
	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(cfg.StoreInterval)
		go func() {
			for {
				select {
				case <-ticker.C:
					if err := storage.SaveToFile(cfg.FileStoragePath); err != nil {
						log.Error().Err(err).Msg("failed to save metrics to file")
					} else {
						log.Info().Msg("metrics saved to file")
					}
				case <-done:
					ticker.Stop()
					return
				}
			}
		}()
	}

	// 3) Запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	// 4) Graceful shutdown по сигналу
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Info().Msg("Shutting down server...")

	// Останавливаем периодическую запись
	close(done)

	// Сохраняем один последний раз
	if err := storage.SaveToFile(cfg.FileStoragePath); err != nil {
		log.Error().Err(err).Msg("failed to save metrics on shutdown")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	log.Info().Msg("Server stopped")
}
