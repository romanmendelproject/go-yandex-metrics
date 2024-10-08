package report

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
)

func handlerServer(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

func TestSendMetrics(t *testing.T) {
	config.ParseFlags()
	type args struct {
		name   string
		metric metrics.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Good Gauge Test",
			args:    args{"TestGauge", metrics.Metric{ID: "test", MType: "gauge", Value: utils.GetFloatPtr(float64(0.5))}},
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
		data := []metrics.Metric{
			{ID: "Alloc", MType: "gauge", Value: utils.GetFloatPtr(1)},
		}
		jsonValue, _ := json.Marshal(data)
		t.Run(tt.name, func(t *testing.T) {

			if err := sendMetric(jsonValue, "http://127.0.0.1:8080/updates/"); (err != nil) != tt.wantErr {
				t.Errorf("reportMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReportSingleMetric(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	metricsChannel := make(chan *[]metrics.Metric, 1)

	wg.Add(1)
	go ReportSingleMetric(ctx, &wg, metricsChannel)

	metric := metrics.Metric{
		ID:    "metric1",
		MType: "gauge",
		Value: float64Ptr(1.1),
	}

	metricsChannel <- &[]metrics.Metric{metric}

}

func TestReportBatchMetric(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	metricsChannel := make(chan *[]metrics.Metric, 1)

	wg.Add(1)
	go ReportBatchMetric(ctx, &wg, metricsChannel)

	metrics := []metrics.Metric{
		{ID: "metric1", MType: "gauge", Value: float64Ptr(1.1)},
		{ID: "metric2", MType: "counter", Delta: int64Ptr(2)},
	}

	metricsChannel <- &metrics

}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
