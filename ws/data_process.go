package ws_handle

import (
	"github.com/gorilla/websocket"
	log "mdocker/logger"
	"unicode/utf8"
)

/**
  @Author   : bob
  @Datetime : 2023-05-14 下午 12:29
  @File     : data_process.go
  @Desc     :
*/

func ReceiveFromClient(cli *clientStruct) {
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

func SendFromServer(cli *clientStruct) {
	for msgByte := range cli.sendChan {
		err := cli.conn.WriteMessage(websocket.TextMessage, msgByte)
		if err != nil {
			log.Log.Error("fal to send data to client, ", err)
			break
		}
		log.Log.Info(string(msgByte))
	}
}

func ReceiveFromDocker(cli *clientStruct) error {
	for {
		buffer := make([]byte, 128)
		n, err := cli.reader.Read(buffer)
		if !utf8.Valid(buffer) {
			log.Log.Errorf("Received message from client contains invalid UTF-8: %v", buffer)
			continue
		}
		if err != nil {
			log.Log.Error("fail to read container logs, ", err)
			return err
		}
		if n > 0 {
			cli.sendChan <- buffer
		}
	}
}
