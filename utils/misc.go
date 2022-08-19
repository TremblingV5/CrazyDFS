package utils

import (
	"net"
	"os"
	"strings"
)

func GetConfigPath() string {
	execFilePath := os.Args[0]

	if strings.Contains(execFilePath, ".test") {
		return "../config/total.yaml"
	} else {
		return "config/total.yaml"
	}
}

func GetIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
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
