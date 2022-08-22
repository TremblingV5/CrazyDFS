package utils

import (
	"io/ioutil"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetConfigPath() string {
	execFilePath := os.Args[0]

	if strings.Contains(execFilePath, ".test") {
		return "../config/total.yaml"
	} else {
		return "config/total.yaml"
	}
}

func GetIPDebug() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		WriteLog(
			"error", "error",
			"message", err.Error(),
		)
	}
	defer conn.Close()

	local := conn.LocalAddr().(*net.UDPAddr)
	return local.IP
}

func Update[T comparable, Y any](src map[T]Y, target map[T]Y) map[T]Y {
	for k, v := range target {
		src[k] = v
	}
	return src
}

func IsExist(arr []string, element string) bool {
	for _, item := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func CheckAndMkdir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

func ReadYaml[T any](path string, t *T) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		WriteLog(
			"error", "Read yaml " + path + " defeat",
		)
	}

	if err := yaml.Unmarshal(bytes, t); err != nil {
		WriteLog(
			"error", "Unmarshal " + path + " defeat",
		)
	}
}

func WriteYaml[T any](path string, t *T) {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		WriteLog(
			"error", "Marshal yaml " + path + " defeat",
		)
	}

	if err := ioutil.WriteFile(path, bytes, 0777); err != nil {
		WriteLog(
			"error", "Write yaml file " + path + " defeat",
		)
	}
}
