package handler

import (
	"context"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	ctx    context.Context
	client influxdb2.Client
	once   sync.Once
)

func GetInfluxdbClient() (influxdb2.Client, context.Context) {
	url := "http://localhost:8086"
	token := "nXom-eG4c_SwYgoDsxFOrfz_hIg0_je6tSDTkC9MssBxl5F0eyVzSEIOazx_89513o4Y5Ld2NOwsllCK41L3xg=="
	once.Do(func() {
		client = influxdb2.NewClient(url, token)
		ctx = context.Background()
		log.Infof("influxdb client is initialized success, client: %p", &client)
	})
	log.Infof("Get influxdb client success, client: %p", &client)
	return client, ctx
}
