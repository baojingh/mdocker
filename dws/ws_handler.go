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
	statsChan chan types.StatsJSON
}

func ContainerStats(w http.ResponseWriter, r *http.Request) {
	conn, _ := RegisterWsClient(w, r)
	cli := &clientStruct{
		conn:      conn.conn,
		statsChan: make(chan types.StatsJSON, 100),
	}
	defer func() {
		log.Info("cli.conn is closed")
		cli.conn.Close()
	}()

	// 传递参数
	// ws://127.0.0.1:8081/stats?id=a315b7da073d
	containerId := "79bc6f76f667"

	ctx, reader, err := container.ContainerStats(containerId)
	if err != nil {
		log.Error("fail to get the container stats reader, ", err)
		return
	}
	defer func() {
		log.Info("stats reader is closed")
		reader.Close()
	}()

	decoder := json.NewDecoder(reader)
	var statsValue types.StatsJSON

	go func() {
		for {
			// 可以同时监听多个通道的数据流动，通过 case 分支来处理具体的通道操作。
			// 当有多个通道同时可用时，select 会随机选择一个可用的通道进行操作。
			// 当所有的通道都阻塞时，select 可以执行默认的 default 分支，实现非阻塞的操作。
			select {
			case <-ctx.Done():
				reader.Close()
				log.Warn("Stop logging metrics")
				return
			default:
				if err := decoder.Decode(&statsValue); err == io.EOF {
					log.Warn("Receive EOF flag")
					return
				} else if err != nil {
					log.Error("Something Error occured", err)
					return
				} else {
					cli.statsChan <- statsValue
				}
			}
		}
	}()

	// 当通道为空时，range 会阻塞等待，直到通道中有数据或者通道被关闭。
	for statJSON := range cli.statsChan {
		statsJSONBytes, err := json.MarshalIndent(statJSON, "", "  ")
		if err != nil {
			panic(err)
		}
		log.Info(string(statsJSONBytes))
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
