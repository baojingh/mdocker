package config

/**
  @Author   : bob
  @Datetime : 2023-05-19 下午 11:16
  @File     : utils.go
  @Desc     :
*/

/**
docker:
  host: 192.168.1.130
  port: 2375
  protocal: tcp

websocket:
  port: 8081

*/

import (
	"github.com/spf13/viper"
	log "mdocker/logger"
)

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
		log.Log.Error("fail to read config file, ", err)
	}
	if err := viper.Unmarshal(&MDocker); err != nil {
		log.Log.Error("fail to get config object, ", err)
	}
}
