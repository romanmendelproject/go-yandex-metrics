package handlers

import (
	"context"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
	pb "github.com/romanmendelproject/go-yandex-metrics/proto"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
	log "github.com/sirupsen/logrus"
)

// ProtoServiceHandlers data for gRPC server
type ProtoServiceHandlers struct {
	storage Storage
}

// NewProtoHandlers создает объект обработчика запросов
func NewProtoHandlers(storage Storage) *ProtoServiceHandlers {
	return &ProtoServiceHandlers{
		storage: storage,
	}
}

// ValueGauge имплементирует ValueGauge
func (h *ProtoServiceHandlers) ValueGauge(ctx context.Context, in *pb.ValueGaugeRequest) (*pb.ValueGaugeResponse, error) {
	value, err := h.storage.GetGauge(context.Background(), in.ID)
	if err != nil {
		return nil, err
	}

	var response pb.ValueGaugeResponse

	log.Info("gRPC ValueGauge:", "_")

	response.Value = value

	log.Infof("Received ID: %s, Value: %s", in.ID, value)

	return &response, nil
}

// ValueCounter имплементирует ValueCounter
func (h *ProtoServiceHandlers) ValueCounter(ctx context.Context, in *pb.ValueCounterRequest) (*pb.ValueCounterResponse, error) {
	value, err := h.storage.GetCounter(context.Background(), in.ID)
	if err != nil {
		return nil, err
	}

	var response pb.ValueCounterResponse

	log.Info("gRPC ValueCounter:", "_")

	response.Delta = value

	log.Infof("Received ID: %s, Value: %s", in.ID, value)

	return &response, nil
}

// ProtoHandler  UpdateBatch implements UpdateBatch
func (h *ProtoServiceHandlers) UpdateBatch(ctx context.Context, in *pb.UpdateBatchRequest) (*pb.UpdateBatchResponse, error) {
	var response pb.UpdateBatchResponse
	ms := []metrics.Metric{}

	for _, metric := range in.Metric {
		ms = append(ms, metrics.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: utils.ToPointer(metric.Delta),
			Value: utils.ToPointer(metric.Value),
		})
	}

	err := h.storage.SetBatch(ctx, ms)
	if err != nil {
		log.Error("gRPC UpdateMetric, s.DBPostgres.SetBatchMetrics:", "about ERR"+err.Error())
	}

	return &response, nil
}
