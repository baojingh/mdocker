package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	handler "mdocker/handler"
	logger "mdocker/logger"
	"os"
	"os/signal"
	"syscall"
)

var log = logger.New()

func main() {

	// statsChan := make(chan types.StatsJSON)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		sig := <-sigChan
		log.Warnf("Receive signal %s and exit", sig)
		os.Exit(0)
	}()

	// handler.StatsProducer("vcfa", statsChan)
	containerList := handler.ContainerList()
	for _, ele := range containerList {
		log.Infof("id: %s, name: %s", ele.Id, ele.Name)
	}

	log.Info("mdocker service starts success")
}
