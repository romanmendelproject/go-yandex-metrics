package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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
	var url_params URLParams

	url_data := strings.Split(url[1:], "/")
	if len(url_data) != 4 {
		return url_params, errors.New("Error parameters coumt from URL")
	}
	url_params.metricType = url_data[1]
	url_params.metricName = url_data[2]
	url_params.metricValue = url_data[3]

	return url_params, nil
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	url_params, err := parseURL(req.URL.Path)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
	}

	switch url_params.metricType {
	case "counter":
		valueInt, err := strconv.ParseInt(url_params.metricValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

		}

		if storage.counter == nil {
			storage.counter = make(map[string]int64)
		}
		storage.counter[url_params.metricName] += valueInt
		res.WriteHeader(http.StatusOK)
	case "gauge":
		valueFloat, err := strconv.ParseFloat(strings.TrimSpace(url_params.metricValue), 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}

		if storage.gauge == nil {
			storage.gauge = make(map[string]float64)
		}
		storage.gauge[url_params.metricName] = valueFloat
		res.WriteHeader(http.StatusOK)
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
