package agent

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		err := s.sendMetric("gauge", name, strconv.FormatFloat(value, 'f', -1, 64))
		if err != nil {
			return err
		}
	}

	for name, value := range counters {
		err := s.sendMetric("counter", name, strconv.FormatInt(value, 10))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sender) sendMetric(metricType, name, value string) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", s.serverAddr, metricType, name, value)

	start := time.Now()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("metric_type", metricType).
			Str("name", name).
			Str("value", value).
			Msg("failed to create HTTP request")
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Error().
			Err(err).
			Str("url", url).
			Str("metric_type", metricType).
			Str("name", name).
			Str("value", value).
			Dur("duration", duration).
			Msg("HTTP request failed")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("status", resp.StatusCode).
			Str("url", url).
			Str("metric_type", metricType).
			Str("name", name).
			Str("value", value).
			Dur("duration", duration).
			Msg("unexpected HTTP response status")
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	log.Info().
		Str("method", req.Method).
		Str("url", url).
		Int("status", resp.StatusCode).
		Str("metric_type", metricType).
		Str("name", name).
		Str("value", value).
		Dur("duration", duration).
		Msg("sent metric successfully")

	return nil
}
