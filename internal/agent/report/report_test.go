package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
)

func handlerServer(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func Test_ReportMetrics(t *testing.T) {
	config.ParseFlags()
	type args struct {
		name   string
		metric metrics.MetricGauge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good Gauge Test",
			args:    args{"TestGauge", metrics.MetricGauge{Type: "gauge", Value: float64(0.5)}},
			wantErr: false,
		},
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(handlerServer))
	if err = server.Listener.Close(); err != nil {
		t.Fatal("failed to close default listener:", err)
	}
	server.Listener = listener

	server.Start()
	defer server.Close()
	for _, tt := range tests {
		var metrics metrics.Metrics
		metrics.Init()
		metrics.DataGauge[tt.args.name] = tt.args.metric

		t.Run(tt.name, func(t *testing.T) {
			if err := ReportMetrics(&metrics); (err != nil) != tt.wantErr {
				t.Errorf("reportMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}