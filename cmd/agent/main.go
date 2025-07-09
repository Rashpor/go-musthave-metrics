package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := agent.NewConfig()

	log.Printf("Agent config: addr=%s, report=%ds, poll=%ds", cfg.Address, cfg.ReportInterval, cfg.PollInterval)

	collector := agent.NewCollector()
	sender := agent.NewSender("http://" + cfg.Address)

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			collector.Collect()
		case <-reportTicker.C:
			gauges, counters := collector.GetMetrics()
			if err := sender.Send(gauges, counters); err != nil {
				log.Printf("failed to send metrics: %v", err)
			}
		}
	}
}
