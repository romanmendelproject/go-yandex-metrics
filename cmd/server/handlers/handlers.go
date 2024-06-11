package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"
	log "github.com/sirupsen/logrus"
)

func HandleBadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func HandleStatusNotFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}

func UpdateGauge(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		storage.SetGauge(urlParams.metricName, valueFloat)
		res.WriteHeader(http.StatusOK)
	}
}

func UpdateCounter(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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

		storage.SetCounter(urlParams.metricName, valueInt)

		res.WriteHeader(http.StatusOK)
	}
}

func ValueGauge(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		urlParams, err := ParseURLValue(req.URL.Path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		value, err := storage.GetGauge(urlParams.metricName)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(res, fmt.Sprintf("%f", value))
	}
}

func ValueCounter(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		urlParams, err := ParseURLValue(req.URL.Path)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		value, err := storage.GetCounter(urlParams.metricName)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(res, fmt.Sprintf("%d", value))
	}
}

func AllData(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		values := storage.GetAll()

		for i, value := range values {
			io.WriteString(res, fmt.Sprintf("%d type = %s  name = %s value = %v", i, value.Type, value.Name, value.Value))
		}
	}
}
