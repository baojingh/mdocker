package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	logger "mdocker/logger"
	ws_handle "mdocker/ws"
)

var log = logger.New()

func main() {
	log.Info("mdocker service starts")
	ws_handle.StartWs()
}
