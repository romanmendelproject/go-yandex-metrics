package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

type URLParams struct {
	metricType  string
	metricName  string
	metricValue string
}

var storage MemStorage

func parseURL(url string) (URLParams, error) {
	var urlParams URLParams

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 4 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.metricType = urlData[1]
	urlParams.metricName = urlData[2]
	urlParams.metricValue = urlData[3]

	return urlParams, nil
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	urlParams, err := parseURL(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
	}

	switch urlParams.metricType {
	case "counter":
		valueInt, err := strconv.ParseInt(urlParams.metricValue, 10, 64)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusBadRequest)

		}

		if storage.counter == nil {
			storage.counter = make(map[string]int64)
		}
		storage.counter[urlParams.metricName] += valueInt
		res.WriteHeader(http.StatusOK)
	case "gauge":
		valueFloat, err := strconv.ParseFloat(strings.TrimSpace(urlParams.metricValue), 64)
		if err != nil {
			log.Error(err)
			res.WriteHeader(http.StatusBadRequest)
		}

		if storage.gauge == nil {
			storage.gauge = make(map[string]float64)
		}
		storage.gauge[urlParams.metricName] = valueFloat
		res.WriteHeader(http.StatusOK)
	default:
		res.WriteHeader(http.StatusBadRequest)
		log.Error("Error type of metric")
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
