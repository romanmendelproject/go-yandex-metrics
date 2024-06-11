package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	activateEnvFlags()
}

func activateEnvFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
}
