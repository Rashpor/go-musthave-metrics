package main

import (
	"flag"
	"log"
	"time"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
)

func main() {
	addr := flag.String("a", "localhost:8080", "address of the metrics server")
	reportInterval := flag.Int("r", 10, "report interval in seconds")
	pollInterval := flag.Int("p", 2, "poll interval in seconds")
	flag.Parse()

	collector := agent.NewCollector()
	sender := agent.NewSender("http://" + *addr)

	pollTicker := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(*reportInterval) * time.Second)
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
