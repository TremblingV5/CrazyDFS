package utils

import (
	"io/ioutil"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"gopkg.in/yaml.v3"
)

var TotalConf *items.Total
var Conf any

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

func InitNodeConfig[T items.DN | items.NN | items.Client](config T, path string) (T, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return config, err
	}

	Conf = &config
	return config, nil
}
