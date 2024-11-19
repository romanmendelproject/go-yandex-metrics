package network

import (
	"net/http"

	"github.com/romanmendelproject/go-yandex-metrics/utils"
)

func XrealIPMiddleware(trustedsubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(res http.ResponseWriter, req *http.Request) {
			xrealip := req.Header.Get("X-Real-IP")
			if trustedsubnet != "" {
				result := utils.ISinTrustedNetwork(xrealip, trustedsubnet)
				if !result {
					res.WriteHeader(http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(logFn)
	}
}
