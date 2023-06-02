package main

import (
	"net/http"
	"test/wsim2/server/gateway/controller"
	"test/wsim2/server/mq/consumer"
	"test/wsim2/server/mq/producer"
)

func main() {
	go producer.Produce()
	go consumer.Consume()
	go controller.Mq2ws()
	http.HandleFunc("/", controller.CreateConn)
	http.ListenAndServe(":1111", nil)
}
