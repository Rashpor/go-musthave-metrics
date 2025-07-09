package server

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address         string
	StoreInterval   time.Duration // в секундах
	FileStoragePath string
	Restore         bool
}

const (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoragePath   = "metrics-db.json"
)

func NewConfig() *Config {
	// старые флаги
	flagAddr := flag.String("a", "", "address to run HTTP server on")

	// новые
	flagStore := flag.Int("i", 0, "interval for storing metrics to disk in seconds (0 = sync on each update)")
	flagFile := flag.String("f", "", "file path for storing metrics")
	flagRestore := flag.Bool("r", false, "restore metrics from file on startup")

	flag.Parse()

	// адрес
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = *flagAddr
	}
	if addr == "" {
		addr = defaultAddress
	}

	// интервал
	var interval time.Duration
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil {
			interval = time.Duration(sec) * time.Second
		}
	} else if *flagStore != 0 {
		interval = time.Duration(*flagStore) * time.Second
	} else {
		interval = defaultStoreInterval
	}

	// путь
	path := os.Getenv("FILE_STORAGE_PATH")
	if path == "" {
		if *flagFile != "" {
			path = *flagFile
		} else {
			path = defaultStoragePath
		}
	}

	// restore
	restore := false
	if v := os.Getenv("RESTORE"); v != "" {
		restore = (v == "true" || v == "1")
	} else {
		restore = *flagRestore
	}

	return &Config{
		Address:         addr,
		StoreInterval:   interval,
		FileStoragePath: path,
		Restore:         restore,
	}
}
