package handlers

import (
	"errors"
	"strings"
)

type URLParams struct {
	metricType, metricName, metricValue string
}

func ParseURLUpdate(url string) (URLParams, error) {
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

func ParseURLValue(url string) (URLParams, error) {
	var urlParams URLParams

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 3 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.metricType = urlData[1]
	urlParams.metricName = urlData[2]

	return urlParams, nil
}
