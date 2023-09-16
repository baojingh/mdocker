package config

/**
  @Author   : bob
  @Datetime : 2023-05-19 下午 11:16
  @File     : utils.go
  @Desc     :
*/

import (
	"mdocker/logger"

	"github.com/spf13/viper"
)

// logger in container package
var log = logger.New()

type DockerInfo struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type WebsocketInfo struct {
	Port string `yaml:"port"`
}

type MDockerInfo struct {
	Docker    DockerInfo    `yaml:"docker"`
	Websocket WebsocketInfo `yaml:"websocket"`
}

var MDocker MDockerInfo

func init() {
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("fail to read config file, ", err)
	}
	if err := viper.Unmarshal(&MDocker); err != nil {
		log.Error("fail to get config object, ", err)
	}
}
