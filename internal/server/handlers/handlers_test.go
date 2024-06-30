package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestServiceHandlers_UpdateGauge(t *testing.T) {
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
	storage := storage.NewMemStorage()

	handler := NewHandlers(storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handler.UpdateGauge(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestServiceHandlers_UpdateCounter(t *testing.T) {
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
	storage := storage.NewMemStorage()

	handler := NewHandlers(storage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, nil)

			w := httptest.NewRecorder()
			handler.UpdateCounter(w, request)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestServiceHandlers_UpdateJSON(t *testing.T) {
	type args struct {
		httpMethod string
		path       string
		body       any
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantValue      metrics.Metric
	}{
		{
			name: "Good update counter",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/",
				body:       metrics.Metric{ID: "test", MType: "gauge", Value: utils.GetFloatPtr(float64(0.5))},
			},
			wantStatusCode: http.StatusOK,
			wantValue:      metrics.Metric{ID: "test", MType: "gauge", Value: utils.GetFloatPtr(float64(0.5))},
		},
		{
			name: "Bad (Incorrect type)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/",
				body:       metrics.Metric{ID: "test"},
			},
			wantStatusCode: http.StatusBadRequest,
			wantValue:      metrics.Metric{},
		},
		{
			name: "Bad (bad http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "/update/",
				body:       metrics.Metric{ID: "test", MType: "gauge", Value: utils.GetFloatPtr(float64(0.5))},
			},
			wantStatusCode: http.StatusBadRequest,
			wantValue:      metrics.Metric{},
		},
		{
			name: "Bad (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/",
				body: map[string]any{
					"id":    "Test3",
					"type":  "counter",
					"delta": "test",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantValue:      metrics.Metric{},
		},
		{
			name: "Bad (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "/update/",
				body: map[string]any{
					"id":    "Test3",
					"type":  "counter",
					"delta": 0.5,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantValue:      metrics.Metric{},
		},
	}

	stor := storage.NewMemStorage()

	handler := NewHandlers(stor)
	var buf bytes.Buffer
	var metric Metric

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, _ := json.Marshal(tt.args.body)

			request := httptest.NewRequest(tt.args.httpMethod, tt.args.path, bytes.NewBuffer(resp))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.UpdateJSON(w, request)

			response := w.Result()
			defer response.Body.Close()

			_, err := buf.ReadFrom(response.Body)
			if err != nil {
				log.Error(err)
			}

			_ = json.Unmarshal(buf.Bytes(), &metric)

			if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
				log.Error(err)
			}
			require.Equal(t, response.StatusCode, tt.wantStatusCode)
			if tt.wantValue != (metrics.Metric{}) {
				require.Equal(t, tt.args.body, tt.wantValue)
			}
		})
	}
}
