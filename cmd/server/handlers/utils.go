package handlers

import (
	"errors"
	"strings"
)

type URLParamsUpdate struct {
	metricType  string
	metricName  string
	metricValue string
}

type URLParamsValue struct {
	metricType string
	metricName string
}

func ParseURLUpdate(url string) (URLParamsUpdate, error) {
	var urlParams URLParamsUpdate

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 4 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.metricType = urlData[1]
	urlParams.metricName = urlData[2]
	urlParams.metricValue = urlData[3]

	return urlParams, nil
}

func ParseURLValue(url string) (URLParamsValue, error) {
	var urlParams URLParamsValue

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 3 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.metricType = urlData[1]
	urlParams.metricName = urlData[2]

	return urlParams, nil
}
