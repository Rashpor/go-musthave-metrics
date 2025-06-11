package main

import (
	"log"
	"time"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
)

func main() {
	const (
		pollInterval   = 2 * time.Second
		reportInterval = 10 * time.Second
		serverAddress  = "http://localhost:8080"
	)

	collector := agent.NewCollector()
	sender := agent.NewSender(serverAddress)

	// Тикер сбора метрик
	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for range ticker.C {
			collector.Collect()
		}
	}()

	// Тикер отправки метрик
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := sender.Send(collector.GetMetrics())
		if err != nil {
			log.Printf("failed to send metrics: %v", err)
		}
	}
}
