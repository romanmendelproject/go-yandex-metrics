package utils

import (
	"errors"
	"strings"
)

type URLParams struct {
	MetricType, MetricName, MetricValue string
}

func ParseURLUpdate(url string) (URLParams, error) {
	var urlParams URLParams

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 4 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.MetricType = urlData[1]
	urlParams.MetricName = urlData[2]
	urlParams.MetricValue = urlData[3]

	return urlParams, nil
}

func ParseURLValue(url string) (URLParams, error) {
	var urlParams URLParams

	urlData := strings.Split(url[1:], "/")
	if len(urlData) != 3 {
		return urlParams, errors.New("error parameters coumt from URL")
	}
	urlParams.MetricType = urlData[1]
	urlParams.MetricName = urlData[2]

	return urlParams, nil
}

func GetFloatPtr(v float64) *float64 {
	return &v
}
