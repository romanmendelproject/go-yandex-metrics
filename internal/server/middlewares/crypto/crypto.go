package crypto

import (
	"bytes"
	"io"
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/internal/crypto"
)

// CryptoMiddleware декодирует запросы к серверу
func CryptoMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(res http.ResponseWriter, req *http.Request) {
			if key != "" {
				// decrypt request body
				data, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(res, err.Error(), http.StatusBadRequest)
					return
				}

				// decrypt only non-empty data
				if len(data) > 0 {
					decryptBody, err := crypto.Decrypt(key, string(data))
					if err != nil {
						http.Error(res, err.Error(), http.StatusBadRequest)
						return
					}
					// возвращаем тело запроса
					req.Body = io.NopCloser(bytes.NewReader([]byte(decryptBody)))
				}
			}

			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(logFn)
	}
}
