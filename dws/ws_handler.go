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
	log.Infof("Starting server on port %s", wsPort)
	err := http.ListenAndServe(wsPort, nil)
	if err != nil {
		log.Error("Failed to start server:", err)
	}
}
