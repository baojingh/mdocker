package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
}

func main() {
	log.Info("hello")

}
