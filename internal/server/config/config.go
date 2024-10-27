package config

import (
	"bytes"
	"encoding/json"
	"os"

	env "github.com/caarlos0/env/v8"
	"github.com/spf13/pflag"
)

type ClientFlags struct {
	FlagRunAddr     string `env:"ADDRESS" json:"address"`
	LogLevel        string `env:"LOG_LEVEL" json:"debug_level"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DBDSN           string `env:"DBDSN" json:"dbdsn"`
	Key             string `env:"KEY" json:"key"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config          string `env:"CONFIG" json:"config"`
}

func ParseFlags() (*ClientFlags, error) {
	flags := new(ClientFlags)
	pflag.StringVarP(&flags.FlagRunAddr, "address", "a", ":8080", "Address and port to run agent")
	pflag.StringVarP(&flags.LogLevel, "LogLevel", "l", "debug", "debug level")
	pflag.IntVarP(&flags.StoreInterval, "StoreInterval", "i", 5, "store interval")
	pflag.StringVarP(&flags.FileStoragePath, "FileStoragePath", "f", "/tmp/metrics-db.json", "storage file path")
	pflag.BoolVarP(&flags.Restore, "Restore", "r", true, "restore data from file")
	pflag.StringVarP(&flags.DBDSN, "DBDSN", "d", "postgres://username:userpassword@localhost:5432/dbname", "db connection")
	pflag.StringVarP(&flags.Key, "Key", "k", "", "hash key")
	pflag.StringVarP(&flags.CryptoKey, "crypto-key", "e", "/home/user/practicum/go-yandex-metrics/certs/private.pem", "crypto-key")

	pflag.Parse()

	if err := env.Parse(flags); err != nil {
		return nil, err
	}

	return flags, nil
}

func ReadConfig(flags *ClientFlags) (*ClientFlags, error) {
	pflag.StringVarP(&flags.Config, "config", "c", "/home/user/practicum/go-yandex-metrics/cmd/server/config.json", "Path to server config file")

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
