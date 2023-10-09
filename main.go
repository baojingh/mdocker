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
	"sync"
	"syscall"

	"github.com/docker/docker/api/types"
)

var log = logger.New()

func main() {
	log.Info("mdocker service starts to work.")

	wg := &sync.WaitGroup{}
	shutdownChan := make(chan int, 1)
	containerList := handler.ContainerList()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		sig := <-sigChan

		// value count in channnel must corespond to the number of goroutine
		for i := 0; i < len(containerList); i++ {
			shutdownChan <- 1
		}
		log.Warnf("Receive signal %s and preapre for exit", sig)
	}()

	for _, ele := range containerList {
		log.Infof("id: %s, name: %s", ele.Id, ele.Name)
		statsChan := make(chan types.StatsJSON)
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.StatsProducer(ele.Id, statsChan, shutdownChan)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.DbConsumer(statsChan, shutdownChan)

		}()
	}
	wg.Wait()
	log.Info("mdocker service shutdown success")

}
