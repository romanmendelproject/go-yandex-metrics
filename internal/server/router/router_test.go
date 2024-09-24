// Модуль управления эндпоинтами
package router

import (
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	var storage handlers.Storage
	handler := handlers.NewHandlers(storage)
	r := NewRouter(handler)

	assert.NotEmpty(t, r)
}
