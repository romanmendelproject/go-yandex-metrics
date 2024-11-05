// Модуль управления эндпоинтами
package router

import (
	"log"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/stretchr/testify/assert"
)

var cfg *config.ClientFlags

func getCfg() {
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

	config.ReadConfig(cfg)
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

}

func TestMain(m *testing.M) {
	getCfg()
}

func TestRouter(t *testing.T) {
	var storage handlers.Storage
	handler := handlers.NewHandlers(storage)
	r := NewRouter(cfg, handler)

	assert.NotEmpty(t, r)
}
