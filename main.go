package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	logger "mdocker/logger"
)

var log = logger.New()

func main() {
	log.Info("hello")
	// ws_handle.StartWs()
}
