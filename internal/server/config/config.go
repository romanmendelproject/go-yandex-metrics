package config

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr     string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DBDSN           string
	Key             string
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&LogLevel, "l", "debug", "debug level")
	flag.IntVar(&StoreInterval, "i", 5, "store interval")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "storage file path")
	flag.BoolVar(&Restore, "r", true, "restore data from file")
	flag.StringVar(&DBDSN, "d", "", "db connection")
	flag.StringVar(&Key, "k", "", "hash key")

	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		StoreInterval, _ = strconv.Atoi(envStoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		Restore = true
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		DBDSN = envDSN
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		Key = envKey
	}
}
