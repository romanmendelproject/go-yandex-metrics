package config

import (
	"flag"
	"os"
)

var FlagRunAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
}
