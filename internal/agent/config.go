package agent

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

const (
	defaultAddress        = "localhost:8080"
	defaultReportInterval = 10
	defaultPollInterval   = 2
)

func NewConfig() *Config {
	flagAddr := flag.String("a", "", "address of the metrics server")
	flagReport := flag.Int("r", 0, "report interval in seconds")
	flagPoll := flag.Int("p", 0, "poll interval in seconds")
	flag.Parse()

	return &Config{
		Address:        getString("ADDRESS", *flagAddr, defaultAddress),
		ReportInterval: getInt("REPORT_INTERVAL", *flagReport, defaultReportInterval),
		PollInterval:   getInt("POLL_INTERVAL", *flagPoll, defaultPollInterval),
	}
}

func getString(envName, flagVal, fallback string) string {
	if env := os.Getenv(envName); env != "" {
		return env
	}
	if flagVal != "" {
		return flagVal
	}
	return fallback
}

func getInt(envName string, flagVal, fallback int) int {
	if env := os.Getenv(envName); env != "" {
		if parsed, err := strconv.Atoi(env); err == nil {
			return parsed
		}
	}
	if flagVal != 0 {
		return flagVal
	}
	return fallback
}
