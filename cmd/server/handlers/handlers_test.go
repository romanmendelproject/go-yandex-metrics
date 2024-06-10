package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateGauge(t *testing.T) {
	type args struct {
		httpMethod string
		path       string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Good update gauge",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/test/0.1",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Bad (not metric name and value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/gauge/",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Bad (bad http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "/update/gauge/test/0.1",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Bad (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/test",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	storage := storage.InitMemStorage()
	handlerUpdateGauge := UpdateGauge(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handlerUpdateGauge(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	type args struct {
		httpMethod string
		path       string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Good update counter",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/1",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Bad (not metric name and value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Bad (bad http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "/update/counter/test/1",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Bad (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/test",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Bad (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/counter/test/0.15",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	storage := storage.InitMemStorage()
	handlerUpdateCounter := UpdateCounter(&storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handlerUpdateCounter(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}
