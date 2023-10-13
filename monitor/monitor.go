package monitor

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

var datas []string

func Add(str string) string {
	d := []byte(str)
	sd := string(d)
	datas = append(datas, sd)

	return sd
}

// monitor go service and profiling
func MonitorService() {
	go func() {
		for {
			Add("https://github.com/EDDYCJY")
			time.Sleep(1 * time.Second)
		}
	}()
	http.ListenAndServe(":9999", nil)
}
