package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerServer(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func Test_sendMetric(t *testing.T) {
	parseFlags()
	type args struct {
		name   string
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good Gauge Test",
			args:    args{"TestGauge", Metric{Type: "gauge", Value: float64(0.5)}},
			wantErr: false,
		},
		{
			name:    "Bad Gauge Test",
			args:    args{"TestGauge", Metric{Type: "gauge", Value: int64(1)}},
			wantErr: true,
		},
		{
			name:    "Good Counter Test",
			args:    args{"TestCounter", Metric{Type: "counter", Value: int64(1)}},
			wantErr: false,
		},
		{
			name:    "Bad Counter Test",
			args:    args{"TestCounter", Metric{Type: "counter", Value: float64(0.5)}},
			wantErr: true,
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
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMetric(tt.args.name, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
