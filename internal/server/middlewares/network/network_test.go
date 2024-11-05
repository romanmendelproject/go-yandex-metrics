package network

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for XrealIPMiddleware
func TestXrealIPMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		xRealIP        string
		trustedSubnet  string
		expectedStatus int
	}{
		{
			name:           "Allowed IP",
			xRealIP:        "192.168.1.1",
			trustedSubnet:  "192.168.1.0/24",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Forbidden IP",
			xRealIP:        "192.168.2.1",
			trustedSubnet:  "192.168.1.0/24",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No Trusted Subnet",
			xRealIP:        "192.168.2.1",
			trustedSubnet:  "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Real-IP", tt.xRealIP)

			rr := httptest.NewRecorder()
			handler := XrealIPMiddleware(tt.trustedSubnet)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK) // Respond with OK if the middleware allows the request
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
