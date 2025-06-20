package main

import (
	"log"
	"time"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
)

func main() {
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
