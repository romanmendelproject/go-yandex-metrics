package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/logger"
)

func NewRouter(handler *handlers.ServiceHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)

	r.Get("/", handler.AllData)
	r.Post("/", handlers.HandleBadRequest)

	// TODO Не разобрался сходу как прокинуть параметры {mname} и т.д. в handlers.ValueGauge
	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{mname}", handler.ValueGauge)
		r.Get("/counter/{mname}", handler.ValueCounter)
		r.Get("/*", handlers.HandleBadRequest)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{mname}/{mvalue}", handler.UpdateGauge)
		r.Post("/counter/{mname}/{mvalue}", handler.UpdateCounter)
		r.Post("/counter/*", handlers.HandleStatusNotFound)
		r.Post("/gauge/*", handlers.HandleStatusNotFound)
		r.Post("/*", handlers.HandleBadRequest)
	})

	return r
}
