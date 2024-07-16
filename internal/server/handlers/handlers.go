package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
	log "github.com/sirupsen/logrus"
)

type Storage interface {
	SetGauge(ctx context.Context, name string, value float64) error
	SetCounter(ctx context.Context, name string, value int64) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) ([]storage.Value, error)
	Ping(ctx context.Context) error
}

type ServiceHandlers struct {
	storage Storage
}

func NewHandlers(storage Storage) *ServiceHandlers {
	return &ServiceHandlers{
		storage: storage,
	}
}

func HandleBadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func HandleStatusNotFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}

func (h *ServiceHandlers) UpdateGauge(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlParams, err := utils.ParseURLUpdate(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	valueFloat, err := strconv.ParseFloat(strings.TrimSpace(urlParams.MetricValue), 64)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.SetGauge(ctx, urlParams.MetricName, valueFloat)
	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) UpdateCounter(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	urlParams, err := utils.ParseURLUpdate(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	valueInt, err := strconv.ParseInt(urlParams.MetricValue, 10, 64)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.SetCounter(req.Context(), urlParams.MetricName, valueInt)

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) ValueGauge(res http.ResponseWriter, req *http.Request) {
	urlParams, err := utils.ParseURLValue(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	value, err := h.storage.GetGauge(req.Context(), urlParams.MetricName)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(res, fmt.Sprintf("%v", strconv.FormatFloat(value, 'f', -1, 64)))
}

func (h *ServiceHandlers) ValueCounter(res http.ResponseWriter, req *http.Request) {
	urlParams, err := utils.ParseURLValue(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	value, err := h.storage.GetCounter(req.Context(), urlParams.MetricName)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(res, fmt.Sprintf("%d", value))
}

func (h *ServiceHandlers) ValueJSON(res http.ResponseWriter, req *http.Request) {
	var metric, metricResponse storage.Metric
	var buf bytes.Buffer
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		log.Error("incorrect Content-Type")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		log.Error("incorrect id data")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	switch metric.MType {
	case "gauge":
		value, err := h.storage.GetGauge(req.Context(), metric.ID)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metricResponse = storage.Metric{
			ID:    metric.ID,
			MType: "gauge",
			Value: &value,
		}

	case "counter":
		value, err := h.storage.GetCounter(req.Context(), metric.ID)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		metricResponse = storage.Metric{
			ID:    metric.ID,
			MType: "counter",
			Delta: &value,
		}
	default:
		log.Error("incorrect type data")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metricResponse)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServiceHandlers) AllData(res http.ResponseWriter, req *http.Request) {
	values, err := h.storage.GetAll(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	for i, value := range values {
		io.WriteString(res, fmt.Sprintf("%d type = %s  name = %s value = %v", i, value.Type, value.Name, value.Value))
	}
}

func (h *ServiceHandlers) Ping(res http.ResponseWriter, req *http.Request) {
	err := h.storage.Ping(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) UpdateJSON(res http.ResponseWriter, req *http.Request) {
	var metric storage.Metric
	var buf bytes.Buffer
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		log.Error("incorrect Content-Type")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if metric.ID == "" {
		log.Error("incorrect id data")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		err := h.storage.SetGauge(req.Context(), metric.ID, *metric.Value)
		if err != nil {
			log.Error(err)
			return
		}
	case "counter":
		h.storage.SetCounter(req.Context(), metric.ID, *metric.Delta)
		counter, err := h.storage.GetCounter(req.Context(), metric.ID)
		if err != nil {
			log.Error(err)
			return
		}
		metric.Delta = &counter
	default:
		log.Error("incorrect type data")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
