package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/romanmendelproject/go-yandex-metrics/cmd/server/storage"
	log "github.com/sirupsen/logrus"
)

func HandleBadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func UpdateGauge(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Error("incorrect http method")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		urlParams, err := ParseURL(req.URL.Path)
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

		urlParams, err := ParseURL(req.URL.Path)
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
