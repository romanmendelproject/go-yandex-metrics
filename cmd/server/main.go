package main

import (
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	parseFlags()
	log.SetLevel(log.DebugLevel)
	storage := storage.InitMemStorage()

	r := chi.NewRouter()

	r.Get("/", handlers.AllData(&storage))
	r.Post("/", handlers.HandleBadRequest)

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{mname}", handlers.ValueGauge(&storage))
		r.Get("/counter/{mname}", handlers.ValueCounter(&storage))
		r.Get("/*", handlers.HandleBadRequest)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{mname}/{mvalue}", handlers.UpdateGauge(&storage))
		r.Post("/counter/{mname}/{mvalue}", handlers.UpdateCounter(&storage))
		r.Post("/counter/*", handlers.HandleStatusNotFound)
		r.Post("/gauge/*", handlers.HandleStatusNotFound)
		r.Post("/*", handlers.HandleBadRequest)
	})

	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
