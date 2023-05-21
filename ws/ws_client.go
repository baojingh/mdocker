package ws_handle

import (
	"errors"
	"github.com/gorilla/websocket"
	log "mdocker/logger"
	"net/http"
)

/**
  @Author   : bob
  @Datetime : 2023-05-14 下午 12:23
  @File     : ws_client.go
  @Desc     :
*/

type WsClient struct {
	conn           *websocket.Conn
	name           string
	active         bool
	receiveMsgChan chan []byte
	sendMsgChan    chan []byte
}

var wsClientsMap = make(map[*WsClient]bool)

func RegisterWsClient(w http.ResponseWriter, r *http.Request, name string) (*WsClient, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Log.Error("fail to create ws success, ", err)
		return nil, err
	}
	log.Log.Infof("get ws client name %s", name)
	client := &WsClient{
		conn:           conn,
		name:           name,
		active:         true,
		receiveMsgChan: make(chan []byte, 10),
		sendMsgChan:    make(chan []byte, 10),
	}
	wsClientsMap[client] = true
	return client, nil
}

func UnregisterClient(client *WsClient) (bool, error) {
	if client == nil {
		log.Log.Error("the client that request to unregister is nil")
		return false, errors.New("NilErr")
	}
	if _, isExit := wsClientsMap[client]; isExit != true {
		delete(wsClientsMap, client)
		log.Log.Infof("delete client %s success", client)
		return true, nil
	}
	log.Log.Infof("fail to delete client %s ", client)
	return false, errors.New("NotFoundErr")
}

func GetClients() map[*WsClient]bool {
	return wsClientsMap
}
