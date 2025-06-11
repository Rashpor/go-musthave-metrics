package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Rashpor/go-musthave-metrics/internal/agent"
)

func getEnvOrDefaultInt(envName string, fallback int) int {
	if val := os.Getenv(envName); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvOrDefaultString(envName string, fallback string) string {
	if val := os.Getenv(envName); val != "" {
		return val
	}
	return fallback
}

func main() {

	flagAddr := flag.String("a", "", "address of the metrics server")
	flagReport := flag.Int("r", 0, "report interval in seconds")
	flagPoll := flag.Int("p", 0, "poll interval in seconds")
	flag.Parse()

	defaultAddr := "localhost:8080"
	defaultReport := 10
	defaultPoll := 2

	addr := getEnvOrDefaultString("ADDRESS", *flagAddr)
	if addr == "" {
		addr = defaultAddr
	}

	reportInterval := getEnvOrDefaultInt("REPORT_INTERVAL", *flagReport)
	if reportInterval == 0 {
		reportInterval = defaultReport
	}

	pollInterval := getEnvOrDefaultInt("POLL_INTERVAL", *flagPoll)
	if pollInterval == 0 {
		pollInterval = defaultPoll
	}

	log.Printf("Agent config: addr=%s, report=%ds, poll=%ds", addr, reportInterval, pollInterval)

	collector := agent.NewCollector()
	sender := agent.NewSender("http://" + addr)

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
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
