// Модуль управления эндпоинтами
package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/middlewares/compress"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/middlewares/crypto"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/middlewares/hash"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/middlewares/logger"
)

// NewRouter определяет эндпоинты для сервера
func NewRouter(handler *handlers.ServiceHandlers) *chi.Mux {
	r := chi.NewRouter()
	if config.CryptoKey != "" {
		r.Use(crypto.CryptoMiddleware(config.CryptoKey))
	}
	r.Use(logger.RequestLogger)
	r.Use(compress.GzipMiddleware)

	// r.Get("/debug/pprof", http.HandlerFunc(pprof.Index))

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/*", http.HandlerFunc(pprof.Index))
	})

	r.Get("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Get("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Get("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Get("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	r.Get("/", handler.AllData)
	r.Post("/", handlers.HandleBadRequest)

	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{mname}", handler.ValueGauge)
		r.Get("/counter/{mname}", handler.ValueCounter)
		r.Post("/", handler.ValueJSON)
		r.Get("/*", handlers.HandleBadRequest)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{mname}/{mvalue}", handler.UpdateGauge)
		r.Post("/counter/{mname}/{mvalue}", handler.UpdateCounter)
		r.Post("/counter/*", handlers.HandleStatusNotFound)
		r.Post("/gauge/*", handlers.HandleStatusNotFound)
		r.Post("/", handler.UpdateJSON)
		r.Post("/*", handlers.HandleBadRequest)
	})
	r.Route("/updates", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		if config.Key != "" {
			r.Use(hash.HashMiddleware(config.Key))
		}
		r.Post("/", handler.UpdateBatch)
	})

	r.Get("/ping", handler.Ping)

	return r
}
