package utils

import (
	"errors"
	"net"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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

// GetIP - получаем IP адрес запуска метода
func GetIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error("Error service.GetIP", "Произошла ошибка при получении интерфейсов: "+err.Error())
		return ""
	}
	resIP := ""
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Error("Error service.GetIP", "Произошла ошибка при получении адресов интерфейса:"+err.Error())
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					resIP = ipNet.IP.String()
				}
			}
		}
	}
	return resIP
}

// ISinTrustedNetwork - проверяем находится ли IP адрес в диапазоне доверенной сети
func ISinTrustedNetwork(checkIP, cidr string) bool {
	mask := strings.Split(cidr, "/")
	subnetBIT := StringToInt(mask[1])
	if mask[0] != "" && subnetBIT == 0 {
		return false
	}
	ip := net.ParseIP(checkIP)
	ipNet := net.IPNet{
		IP:   net.ParseIP(mask[0]),
		Mask: net.CIDRMask(subnetBIT, 32),
	}
	if ipNet.Contains(ip) {
		log.Info("service.ISinTrustedNetwork", "IP-адрес находится в подсети CIDR")
	} else {
		log.Info("service.ISinTrustedNetwork", "IP-адрес НЕ находится в подсети CIDR")
	}
	return ipNet.Contains(ip)
}

// StringToInt
func StringToInt(strVar string) int {
	intVar, err := strconv.Atoi(strVar)
	if err != nil {
		log.Error("Error StringToInt", "convert; about err: "+err.Error()+"string: "+string(strVar))
		return 0
	}
	return intVar
}
