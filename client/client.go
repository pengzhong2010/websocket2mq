package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"test/wsim2/msg"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	cch := make(chan int)
	u := url.URL{Scheme: "ws", Host: "localhost:1111", Path: "/"}
	log.Println(u.String())
	for i := 0; i < 2; i++ {
		go func(id int) {
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal(err)
				return
			}
			go func(c *websocket.Conn) {
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("recv:%s", message)
				}

			}(conn)
			defer conn.Close()
			// login
			loginData := msg.MsgData{
				Id:     int64(id),
				Action: "signIn",
				Key:    "",
				Value:  "",
			}
			loginMsg, err := json.Marshal(loginData)
			if err != nil {
				log.Println(err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, loginMsg)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("send:%s", loginMsg)

			// data
			ticker := time.NewTicker(time.Second * 1)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					sendData := msg.MsgData{
						Id:     int64(id),
						Action: "data",
						Key:    "",
						Value:  fmt.Sprintf("clientId:%d,ticker:%s", id, t.String()),
					}
					sendMsg, err := json.Marshal(sendData)
					if err != nil {
						log.Println(err)
						return
					}
					//sendMsg := fmt.Sprintf("clientId:%d,ticker:%s", id, t.String())
					err = conn.WriteMessage(websocket.TextMessage, sendMsg)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("send:%s", sendMsg)
				case <-cch:
					log.Println("interrupt")
					// login
					signOutData := msg.MsgData{
						Id:     int64(id),
						Action: "signOut",
						Key:    "",
						Value:  "",
					}
					signOutMsg, err := json.Marshal(signOutData)
					if err != nil {
						log.Println(err)
						return
					}
					err = conn.WriteMessage(websocket.TextMessage, signOutMsg)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("send:%s", signOutMsg)

					err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println(err)
						return
					}
					return

				}
			}

		}(i)
	}
	time.Sleep(time.Second * 3)
	for i := 0; i < 2; i++ {
		cch <- i
	}
	close(cch)
}
