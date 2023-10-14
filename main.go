package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	"mdocker/handler"
	logger "mdocker/logger"
	"net/http"

	// part 1/2 for pprof monitoring
	_ "net/http/pprof"

	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"sync"
	"syscall"

	"github.com/docker/docker/api/types"
)

var log = logger.New()

func main() {

	// monitorMem()
	// monitorCPU()

	log.Info("mdocker service starts to work.")
	wg := &sync.WaitGroup{}
	shutdownChan := make(chan int, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		sig := <-sigChan
		shutdownChan <- 1
		log.Warnf("Receive signal %s and preapre for exit", sig)
	}()

	statsChan := make(chan types.StatsJSON)
	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.StatsProducer(statsChan, shutdownChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.DbConsumer(statsChan)
	}()

	// part 2/2, for program monitoring with pprof
	go func() {
		http.ListenAndServe(":9988", nil)
	}()

	wg.Wait()
	log.Info("mdocker service shutdown success")
}

func monitorMem() {
	f, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func monitorCPU() {

	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}
