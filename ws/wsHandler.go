package wshandle

import (
	"github.com/gorilla/websocket"
	"io"
	"mdocker/container"
	log "mdocker/logger"
	"net/http"
)

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 10:34 AM
  @Description:
*/

var host = ":8081"

//var containerId = "50a866b74a24"

// 客户端结构
type clientStruct struct {
	conn     *websocket.Conn
	reader   io.ReadCloser
	sendChan chan []byte
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

func StartWs() {
	http.HandleFunc("/logs", ContainerLogs)
	log.Log.Infof("Starting server on port %s", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Log.Error("Failed to start server:", err)
	}
}
