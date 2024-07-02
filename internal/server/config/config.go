package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr     string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&LogLevel, "l", "debug", "debug level")
	flag.IntVar(&StoreInterval, "i", 5, "store interval")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "storage file path")
	flag.BoolVar(&Restore, "r", true, "restore data from file")

	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	// if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
	// 	StoreInterval = int64(envStoreInterval)
	// }
	// if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
	// 	FileStoragePath = envFileStoragePath
	// }
	// if envRestore := os.Getenv("RESTORE"); envRestore != "" {
	// 	Restore = envRestore
	// }
}
