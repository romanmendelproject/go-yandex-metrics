package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	log "github.com/sirupsen/logrus"
)

type Metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Storage interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAll() []storage.Value
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
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlParams, err := ParseURLUpdate(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	valueFloat, err := strconv.ParseFloat(strings.TrimSpace(urlParams.metricValue), 64)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.SetGauge(urlParams.metricName, valueFloat)
	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) UpdateCounter(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	urlParams, err := ParseURLUpdate(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	valueInt, err := strconv.ParseInt(urlParams.metricValue, 10, 64)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.SetCounter(urlParams.metricName, valueInt)

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) ValueGauge(res http.ResponseWriter, req *http.Request) {
	urlParams, err := ParseURLValue(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	value, err := h.storage.GetGauge(urlParams.metricName)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(res, fmt.Sprintf("%v", strconv.FormatFloat(value, 'f', -1, 64)))
}

func (h *ServiceHandlers) ValueCounter(res http.ResponseWriter, req *http.Request) {
	urlParams, err := ParseURLValue(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	value, err := h.storage.GetCounter(urlParams.metricName)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(res, fmt.Sprintf("%d", value))
}

func (h *ServiceHandlers) ValueJSON(res http.ResponseWriter, req *http.Request) {
	var metric, metricResponse Metric
	var buf bytes.Buffer
	if req.Method != http.MethodPost {
		log.Error("incorrect http method")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if res.Header().Get("Content-Type") != "application/json" {
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
		value, err := h.storage.GetGauge(metric.ID)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		metricResponse = Metric{
			ID:    metric.ID,
			MType: "gauge",
			Value: value,
		}
	case "counter":
		value, err := h.storage.GetCounter(metric.ID)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		metricResponse = Metric{
			ID:    metric.ID,
			MType: "counter",
			Delta: value,
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
	values := h.storage.GetAll()

	for i, value := range values {
		io.WriteString(res, fmt.Sprintf("%d type = %s  name = %s value = %v", i, value.Type, value.Name, value.Value))
	}
}

func (h *ServiceHandlers) UpdateJSON(res http.ResponseWriter, req *http.Request) {
	var metric Metric
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
		h.storage.SetGauge(metric.ID, metric.Value)
	case "counter":
		h.storage.SetCounter(metric.ID, metric.Delta)
		counter, err := h.storage.GetCounter(metric.ID)
		if err != nil {
			log.Error(err)
			return
		}
		metric.Delta = counter
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
