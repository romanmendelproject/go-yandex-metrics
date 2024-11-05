package config

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/spf13/pflag"
)

type ClientFlags struct {
	FlagReqAddr          string `env:"ADDRESS" json:"address"`
	ReportSingleInterval int    `env:"REPORT_INTERVAL" json:"report_single_interval"`
	ReportBatchInterval  int    `env:"REPORT_BATCH_INTERVAL" json:"report_batch_interval"`
	PollInterval         int    `env:"POLL_INTERVAL" json:"poll_interval"`
	Key                  string `env:"KEY" json:"key"`
	RateLimit            int    `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKey            string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config               string `env:"CONFIG" json:"config"`
}

func ParseFlags() (*ClientFlags, error) {
	flags := new(ClientFlags)
	pflag.StringVarP(&flags.FlagReqAddr, "address", "a", "localhost:8080", "Address and port to run agent")
	pflag.IntVarP(&flags.ReportSingleInterval, "report-single-interval", "r", 5,
		"Send metrics single method to server")
	pflag.IntVarP(&flags.ReportBatchInterval, "report-batch-interval", "b", 30,
		"Send metrics batch method to server")
	pflag.IntVarP(&flags.PollInterval, "poll-interval", "p", 2,
		"Wait interval in seconds before reading metrics")
	pflag.StringVarP(&flags.Key, "key", "k", "",
		"Hash key to calculate hash sum")
	pflag.IntVarP(&flags.RateLimit, "rateLimit", "l", 2,
		"Max count of parallel outbound requests to server")
	pflag.StringVarP(&flags.CryptoKey, "crypto-key", "e", "./certs/public.pem", "Path to public key RSA to encrypt messages")

	pflag.Parse()

	if err := env.Parse(flags); err != nil {
		return nil, err
	}

	return flags, nil
}

func ReadConfig(flags *ClientFlags) (*ClientFlags, error) {
	pflag.StringVarP(&flags.Config, "config", "c", "./cmd/agent/config.json", "Path to agent config file")

	pflag.Parse()

	data, err := os.ReadFile(flags.Config)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(&flags); err != nil {
		return nil, err
	}
	return flags, nil
}
