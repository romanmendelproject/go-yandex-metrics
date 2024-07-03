package main

import (
	"net/http"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/logger"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/router"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()
	logger.SetLogLevel(config.LogLevel)
	storage := storage.NewMemStorage(config.FileStoragePath)

	handler := handlers.NewHandlers(storage)
	r := router.NewRouter(handler)

	if config.Restore {
		err := storage.RestoreFromFile()
		if err != nil {
			log.Error(err)
		}
	}

	go func() {
		err := http.ListenAndServe(config.FlagRunAddr, r)
		if err != nil {
			panic(err)
		}
	}()
	for {
		time.Sleep(time.Second * time.Duration(config.StoreInterval))
		go func() {
			err := storage.SaveToFile()
			if err != nil {
				log.Error(err)
			}
		}()
	}

}
