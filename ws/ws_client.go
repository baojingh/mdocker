package wshandle

import (
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Log.Error("fail to create ws success, ", err)
		return nil, err
	}
	log.Log.Infof("get ws success")
	return conn, nil
}
