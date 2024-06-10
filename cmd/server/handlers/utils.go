package handlers

import (
	"errors"
	"strings"
)

type URLParams struct {
	metricType  string
	metricName  string
	metricValue string
}

func ParseURL(url string) (URLParams, error) {
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
