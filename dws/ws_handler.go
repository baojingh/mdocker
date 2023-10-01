package dws

import (
	"encoding/json"
	"io"
	"mdocker/config"
	"mdocker/container"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
)

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 10:34 AM
  @Description:
*/

// 客户端结构
type clientStruct struct {
	conn      *websocket.Conn
	reader    io.ReadCloser
	sendChan  chan []byte
	statsChan chan types.StatsJSON
}

func ContainerStats(w http.ResponseWriter, r *http.Request) {
	conn, _ := RegisterWsClient(w, r)
	cli := &clientStruct{
		conn:      conn.conn,
		sendChan:  make(chan []byte, 100),
		statsChan: make(chan types.StatsJSON, 100),
	}
	defer func() {
		log.Info("cli.conn is closed")
		cli.conn.Close()
	}()

	// 传递参数
	// ws://127.0.0.1:8081/stats?id=a315b7da073d
	containerId := "79bc6f76f667"

	reader, err := container.ContainerStats(containerId)
	if err != nil {
		log.Error("fail to get the container stats reader, ", err)
		return
	}
	defer func() {
		log.Info("stats reader is closed")
		reader.Close()
	}()
	go func() {
		var statsValue types.StatsJSON
		for {
			decoder := json.NewDecoder(reader)
			err := decoder.Decode(&statsValue)
			if err != nil {
				log.Error("cannot decode stats data, ", err)
			}
			cli.statsChan <- statsValue
			if err != nil {
				log.Error("fail to read container logs, ", err)
				// TODO 当容器重启之后，reader会被关闭。在实际场景中需要重新获取reader
				// 此处待优化。临时方案是重新获取reader
				reader, _ = container.ContainerStats(containerId)
			}
		}
	}()

	for stats := range cli.statsChan {

		log.Info(stats)
	}
}

func ContainerList(w http.ResponseWriter, r *http.Request) {
	containers, _ := container.ContainerList()
	// 将数组转换为 JSON
	jsonData, err := json.Marshal(containers)
	if err != nil {
		return
	}

	// 设置响应头，指定内容类型为 application/json
	w.Header().Set("Content-Type", "application/json")
	// 发送 JSON 数据作为 HTTP 响应
	w.Write(jsonData)
}

func StartWebsocket() {
	wsPort := config.MDocker.Websocket.Port
	http.HandleFunc("/list", ContainerList)
	http.HandleFunc("/stats", ContainerStats)
	log.Infof("Starting server on port %s", wsPort)
	err := http.ListenAndServe(wsPort, nil)
	if err != nil {
		log.Error("Failed to start server:", err)
	}
}
