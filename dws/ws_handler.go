package dws

import (
	"encoding/json"
	"io"
	"mdocker/config"
	"mdocker/container"
	user "mdocker/user"
	"net/http"
	"unicode/utf8"

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

func ContainerExec(w http.ResponseWriter, r *http.Request) {
	// 传递参数
	// ws://127.0.0.1:8081/exec?id=a315b7da073d
	containerId := r.URL.Query().Get("id")
	clientName := r.URL.Query().Get("name")

	log.Infof("containerid is %s, client name %s", containerId, clientName)

	wsClient, _ := RegisterWsClient(w, r, clientName)
	defer func() {
		log.Info("ws cli conn is closed")
		wsClient.conn.Close()
	}()

	containerHr, err := container.ContainerExec(containerId)
	if err != nil {
		log.Error("fail to get the container exec, ", err)
		return
	}

	// 获取容器返回的消息，并将其转发到websocket
	go func() {
		buf := make([]byte, 512)
		for {
			nr, err := containerHr.Conn.Read(buf)
			log.Infof("read data from container, %s", string(buf))
			if nr > 0 {
				err := wsClient.conn.WriteMessage(websocket.TextMessage, buf[0:nr])
				if err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	// 读取界面的输入数据，并将其转发到容器中
	for {
		_, message, err := wsClient.conn.ReadMessage()
		if !utf8.Valid(message) {
			log.Errorf("Received message from client contains invalid UTF-8: %v", message)
			continue
		}
		containerHr.Conn.Write(message)
		log.Infof("read message from client, %s", string(message))
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Info("Client close conn success")
			} else {
				log.Error("conn has something wrong, ", err)
			}
			break
		}
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
	log.Info(containers)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonData)
}

func StartWebsocket() {
	wsPort := config.MDocker.Websocket.Port
	// POST
	http.HandleFunc("/login", user.LoginHandler)

	http.HandleFunc("/logout", user.LogoutHandler)
	http.HandleFunc("/secret", user.SecretHandler)
	http.HandleFunc("/exec", ContainerExec)
	http.HandleFunc("/list", ContainerList)
	log.Infof("Starting server on port %s", wsPort)
	err := http.ListenAndServe(wsPort, nil)
	if err != nil {
		log.Error("Failed to start server:", err)
	}
}
