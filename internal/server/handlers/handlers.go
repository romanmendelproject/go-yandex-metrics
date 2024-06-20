package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	log "github.com/sirupsen/logrus"
)

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
	fmt.Println(h.storage)
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

func (h *ServiceHandlers) AllData(res http.ResponseWriter, req *http.Request) {
	values := h.storage.GetAll()

	for i, value := range values {
		io.WriteString(res, fmt.Sprintf("%d type = %s  name = %s value = %v", i, value.Type, value.Name, value.Value))
	}
}
