package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rashpor/go-musthave-metrics/models"
	"github.com/rs/zerolog/log"
)

type Sender struct {
	serverAddr string
	client     *http.Client
}

func NewSender(serverAddr string) *Sender {
	return &Sender{
		serverAddr: serverAddr,
		client:     &http.Client{},
	}
}

func (s *Sender) Send(gauges map[string]float64, counters map[string]int64) error {
	for name, value := range gauges {
		val := value
		metric := models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &val,
		}
		if err := s.sendJSONMetric(metric); err != nil {
			return err
		}
	}

	for name, value := range counters {
		val := value
		metric := models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &val,
		}
		if err := s.sendJSONMetric(metric); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sender) sendJSONMetric(metric models.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Msg("failed to marshal metric to JSON")
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Сжимаем JSON через gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write gzip: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Создаём HTTP-запрос с сжатым телом
	req, err := http.NewRequest(http.MethodPost, s.serverAddr+"/update/", &buf)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Msg("failed to create HTTP request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	start := time.Now()
	resp, err := s.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Dur("duration", duration).
			Msg("HTTP request failed")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("status", resp.StatusCode).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Dur("duration", duration).
			Msg("unexpected HTTP response status")
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	log.Info().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Int("status", resp.StatusCode).
		Str("id", metric.ID).
		Str("type", metric.MType).
		Dur("duration", duration).
		Msg("sent metric successfully")

	return nil
}

/*
func (s *Sender) sendJSONMetric(metric models.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Msg("failed to marshal metric to JSON")
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.serverAddr+"/update/", bytes.NewReader(data))
	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Msg("failed to create HTTP request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := s.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Error().
			Err(err).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Dur("duration", duration).
			Msg("HTTP request failed")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("status", resp.StatusCode).
			Str("id", metric.ID).
			Str("type", metric.MType).
			Dur("duration", duration).
			Msg("unexpected HTTP response status")
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// успешно отправлено
	log.Info().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Int("status", resp.StatusCode).
		Str("id", metric.ID).
		Str("type", metric.MType).
		Dur("duration", duration).
		Msg("sent metric successfully")

	return nil
}
*/
