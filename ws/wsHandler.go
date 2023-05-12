package wshandle

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"mdocker/container"
	"net/http"
)

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 10:34 AM
  @Description:
*/

var host = ":8081"

// 消息体
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端结构
type clientStruct struct {
	conn     *websocket.Conn
	sendChan chan wsMessage
}

// 服务端结构
type serverStruct struct {
	clients        map[*clientStruct]bool
	registerChan   chan *clientStruct
	unregisterChan chan *clientStruct
}

var LogsChan chan []byte
var ExecChan chan []byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2014,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
}

func (c *clientStruct) read() {
	defer func() {
		c.conn.Close()
	}()
	for true {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Error("Server fail to receive data from client, ", err)
			break
		}
		log.Infof("Server receive data from client, %s", string(data))
		msg := wsMessage{
			messageType: messageType,
			data:        data,
		}
		c.sendChan <- msg
	}
}

func (c *clientStruct) write(msgChan chan []byte) {
	defer func() {
		c.conn.Close()
	}()
	for true {
		// 阻塞，直到有数据进入msgChan
		msg := <-msgChan
		log.Infof("get msg from channel, %s", string(msg))
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Info("Failed to write message to WebSocket:", err)
			break
		}
		log.Infof("Server send message: %s", msg)
	}
}

func (s *serverStruct) doHandle(w http.ResponseWriter, r *http.Request, msgChan chan []byte) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("fail to create ws success, ", err)
		return
	}
	log.Infof("get ws success")

	cli := &clientStruct{
		conn:     conn,
		sendChan: make(chan wsMessage, 100),
	}
	log.Info("client has registered success")
	go cli.read()
	go cli.write(msgChan)
}

func ContainerLogs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("fail to create ws success, ", err)
		return
	}
	log.Infof("get ws success")

	cli := &clientStruct{
		conn:     conn,
		sendChan: make(chan wsMessage, 100),
	}
	log.Info("client has registered success")

	go func() {
		reader, _ := container.ConLogs("229dab4f6eaf")
		buffer := make([]byte, 8096)
		for {
			if reader == nil {
				log.Error("reader is nil")
			}
			n, _ := reader.Read(buffer)
			log.Infof("read data size: %d, n is %d", len(buffer), n)
			msgByte := make([]byte, n)
			copy(msgByte, buffer[:n])
			err = cli.conn.WriteMessage(websocket.TextMessage, msgByte)
			if err != nil {
				log.Error("fal to send data to client, ", err)
				continue
			}
			log.Info("server send data to client success, #### ", string(msgByte))
		}
	}()

}

func (s *serverStruct) containerExec(w http.ResponseWriter, r *http.Request) {
	s.doHandle(w, r, ExecChan)
}

func (s *serverStruct) unregister(cli *clientStruct) {
	if _, ok := s.clients[cli]; ok {
		delete(s.clients, cli)
		cli.conn.Close()
	}
}

func (s *serverStruct) register(cli *clientStruct) {
	if cli != nil && cli.conn != nil {
		s.clients[cli] = true
		s.registerChan <- cli
	}
}

func StartWs() {
	LogsChan = make(chan []byte, 100)
	ExecChan = make(chan []byte, 100)
	ser := &serverStruct{
		clients:      make(map[*clientStruct]bool, 100),
		registerChan: make(chan *clientStruct, 100),
	}
	http.HandleFunc("/logs", ContainerLogs)
	http.HandleFunc("/exec", ser.containerExec)
	log.Infof("Starting server on port %s", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Error("Failed to start server:", err)
	}
}
