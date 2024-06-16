package main

import (
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/router"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.ParseFlags()
	log.SetLevel(log.DebugLevel)
	storage := storage.NewMemStorage()

	handler := handlers.NewHandlers(storage)
	r := router.NewRouter(handler)

	err := http.ListenAndServe(config.FlagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
