package main

import (
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	storage := storage.InitMemStorage()

	r := chi.NewRouter()

	r.Get("/", handlers.AllData(&storage))
	r.Post("/", handlers.HandleBadRequest)

	r.Get("/value/gauge/{mname}", handlers.ValueGauge(&storage))
	r.Get("/value/counter/{mname}", handlers.ValueCounter(&storage))

	r.Post("/update/gauge/{mname}/{mvalue}", handlers.UpdateGauge(&storage))
	r.Post("/update/counter/{mname}/{mvalue}", handlers.UpdateCounter(&storage))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
