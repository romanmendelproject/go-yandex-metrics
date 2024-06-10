package main

import (
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	storage := storage.InitMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", handlers.UpdateGauge(&storage))
	mux.HandleFunc("/update/counter/", handlers.UpdateCounter(&storage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
