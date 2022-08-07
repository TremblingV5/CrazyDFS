package utils

import (
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
