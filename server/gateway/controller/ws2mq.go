package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"test/wsim2/msg"
	"test/wsim2/server/mq/consumer"
	"test/wsim2/server/mq/producer"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conns map[int64]*websocket.Conn
)

func init() {
	conns = make(map[int64]*websocket.Conn)
}
func CreateConn(w http.ResponseWriter, r *http.Request) {
	conn, errT := upgrader.Upgrade(w, r, nil)
	if errT != nil {
		log.Println(errT)
		return
	}
	defer conn.Close()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// 关闭连接
			return
		}
		log.Printf("in:%s", string(p))
		// 业务逻辑
		msgData := &msg.MsgData{}
		err = json.Unmarshal(p, msgData)
		if err != nil {
			log.Println(err)
			return
		}
		switch msgData.Action {
		case "signIn":
			signIn(msgData.Id, conn)
		case "signOut":
			signOut(msgData.Id)
		default:
			// mq
			producer.Ch <- string(p)
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte("got it"))
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("out:%s", string(p))
	}
}

func Mq2ws() {
	for {
		select {
		case data := <-consumer.Ch:
			//fmt.Println("consumer:", data)
			msgData := &msg.MsgData{}
			err := json.Unmarshal([]byte(data), msgData)
			if err != nil {
				log.Println(err)
			}
			conn, ok := conns[msgData.Id]
			if !ok {
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
func signIn(id int64, conn *websocket.Conn) {
	conns[id] = conn
}
func signOut(id int64) {
	delete(conns, id)
}
