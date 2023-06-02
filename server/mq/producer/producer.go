package producer

import (
	"log"

	"github.com/Shopify/sarama"
)

var (
	Ch      chan string
	Signals chan string
	ChClose chan string
)

func init() {
	Ch = make(chan string)
	Signals = make(chan string)
	ChClose = make(chan string)
}
func Produce() {
	config := sarama.NewConfig()
	config.Version = sarama.V3_3_1_0
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID = "abc"
	client, err := sarama.NewClient([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatalf("client faild,%s", err.Error())
		return
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalf("producer faild,%s", err.Error())
		return
	}
	defer func() {
		producer.Close()
		client.Close()
	}()
ProducerLoop:
	for {
		select {
		case data := <-Ch:
			msg := &sarama.ProducerMessage{Topic: "ws2mq", Value: sarama.StringEncoder(data)}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("FAILED to send message: %s\n", err)
			} else {
				log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
			}
		case <-Signals:
			break ProducerLoop
		}
	}
	ChClose <- ""
}
