package monitor

import (
	"mdocker/logger"
	"os"
	"runtime/pprof"
)

// logger in container package
var log = logger.New()

func monitorMem() {
	f, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
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
