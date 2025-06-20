package server

import (
	"flag"
	"os"
)

type Config struct {
	Address string
}

const defaultAddress = "localhost:8080"

func NewConfig() *Config {
	flagAddr := flag.String("a", "", "address to run HTTP server on")
	flag.Parse()

	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = *flagAddr
	}
	if addr == "" {
		addr = defaultAddress
	}

	return &Config{
		Address: addr,
	}
}
