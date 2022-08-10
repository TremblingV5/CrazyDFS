package utils

import (
	"io/ioutil"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"gopkg.in/yaml.v3"
)

var TotalConf *items.Total

func InitConfig() error {
	var config items.Total

	configFilePath := GetConfigPath()
	file, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return nil
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return err
	}

	TotalConf = &config
	return nil
}
