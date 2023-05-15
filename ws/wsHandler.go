package wshandle

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
	"io"
	"mdocker/container"
	log "mdocker/logger"
	"net/http"
	"time"
	"unicode/utf8"
)

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 10:34 AM
  @Description:
*/

var host = ":8081"

// 客户端结构
type clientStruct struct {
	conn      *websocket.Conn
	reader    io.ReadCloser
	sendChan  chan []byte
	statsChan chan types.StatsJSON
}

func ContainerLogs(w http.ResponseWriter, r *http.Request) {
	conn, _ := getWs(w, r)
	cli := &clientStruct{
		conn:     conn,
		sendChan: make(chan []byte, 100),
	}
	defer func() {
		log.Log.Info("cli.conn is closed")
		cli.conn.Close()
	}()

	go func() {
		// 传递参数
		// ws://127.0.0.1:8081/logs?id=a315b7da073d
		containerId := r.URL.Query().Get("id")
		for {
			reader, err := container.ContainerLogs(containerId)
			if err != nil {
				log.Log.Error("fail to get the container logs reader, ", err)
				return
			}
			defer func() {
				log.Log.Info("reader is closed")
				reader.Close()
			}()
			cli.reader = reader
			err = ReceiveFromDocker(cli)
			if err != nil {
				log.Log.Error("fail to read container logs, ", err)
				// TODO 当容器重启之后，reader会被关闭。在实际场景中需要重新获取reader
				// 此处待优化。临时方案是重新获取reader
				reader, err = container.ContainerLogs(containerId)
			}
		}
	}()
	go func() {
		ReceiveFromClient(cli)
	}()
	SendFromServer(cli)
}

func ContainerStats(w http.ResponseWriter, r *http.Request) {
	conn, _ := getWs(w, r)
	cli := &clientStruct{
		conn:      conn,
		sendChan:  make(chan []byte, 100),
		statsChan: make(chan types.StatsJSON, 100),
	}
	defer func() {
		log.Log.Info("cli.conn is closed")
		cli.conn.Close()
	}()

	// 传递参数
	// ws://127.0.0.1:8081/stats?id=a315b7da073d
	containerId := r.URL.Query().Get("id")

	reader, err := container.ContainerStats(containerId)
	if err != nil {
		log.Log.Error("fail to get the container stats reader, ", err)
		return
	}
	defer func() {
		log.Log.Info("stats reader is closed")
		reader.Close()
	}()
	go func() {
		var statsValue types.StatsJSON
		for true {
			decoder := json.NewDecoder(reader)
			err := decoder.Decode(&statsValue)
			if err != nil {
				log.Log.Error("cannot decode stats data, ", err)
			}
			cli.statsChan <- statsValue
			if err != nil {
				log.Log.Error("fail to read container logs, ", err)
				// TODO 当容器重启之后，reader会被关闭。在实际场景中需要重新获取reader
				// 此处待优化。临时方案是重新获取reader
				reader, err = container.ContainerStats(containerId)
			}
		}
		time.Sleep(time.Second * 1)
	}()

	go func() {
		ReceiveFromClient(cli)
	}()
	for stats := range cli.statsChan {
		byteArr, err := convertStatsJSONToByte(stats)
		err = cli.conn.WriteMessage(websocket.TextMessage, byteArr)
		if err != nil {
			log.Log.Error("fal to send data to client, ", err)
			break
		}
		log.Log.Info(stats)
	}
}

func convertStatsJSONToByte(stats types.StatsJSON) ([]byte, error) {
	// 使用 json.Marshal() 将 StatsJSON 编码为 JSON 字节数组
	bytes, err := json.Marshal(stats)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ContainerInspect(w http.ResponseWriter, r *http.Request) {
	// 传递参数
	// ws://127.0.0.1:8081/inspect?id=a315b7da073d
	containerId := r.URL.Query().Get("id")
	inspect, err := container.ContainerInspect(containerId)
	if err != nil {
		log.Log.Error("fail to get the container inspect, ", err)
		return
	}
	jsonString, err := json.Marshal(inspect)
	w.Write(jsonString)
}

func ContainerExec(w http.ResponseWriter, r *http.Request) {
	conn, _ := getWs(w, r)
	cli := &clientStruct{
		conn:      conn,
		sendChan:  make(chan []byte, 100),
		statsChan: make(chan types.StatsJSON, 100),
	}
	defer func() {
		log.Log.Info("cli.conn is closed")
		cli.conn.Close()
	}()

	// 传递参数
	// ws://127.0.0.1:8081/exec?id=a315b7da073d
	containerId := r.URL.Query().Get("id")
	log.Log.Infof("get containerid %s", containerId)
	hr, err := container.ContainerExec(containerId)
	if err != nil {
		log.Log.Error("fail to get the container exec, ", err)
		return
	}

	// 转发输入/输出至websocket
	go func() {
		buf := make([]byte, 512)
		for {
			nr, err := hr.Conn.Read(buf)
			log.Log.Infof("read data from container, %s", string(buf))
			if nr > 0 {
				err := conn.WriteMessage(websocket.BinaryMessage, buf[0:nr])
				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	for {
		_, message, err := cli.conn.ReadMessage()
		if !utf8.Valid(message) {
			log.Log.Errorf("Received message from client contains invalid UTF-8: %v", message)
			continue
		}
		log.Log.Infof("read message from client, %s", string(message))
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Log.Info("Client close conn success")
			} else {
				log.Log.Error("conn has something wrong, ", err)
			}
			break
		}
	}
}

func StartWs() {
	//http.HandleFunc("/logs", ContainerLogs)
	//http.HandleFunc("/stats", ContainerStats)
	//http.HandleFunc("/inspect", ContainerInspect)
	http.HandleFunc("/exec", ContainerExec)
	log.Log.Infof("Starting server on port %s", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Log.Error("Failed to start server:", err)
	}
}
