package main

import (
	"fmt"
	"github.com/huhongda/GoToolBox/aliware-kafka"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"os/signal"
)

var jsonConfigString = `
{
  "topics": ["your topic"],
  "servers": ["kafka-ons-internet.aliyun.com:8080"],
  "ak": "Access Key",
  "password": "password",
  "consumerId": "your consumer id",
}
`
var newLogger = log.New(os.Stderr, "", log.LstdFlags)

func main() {
	results := gjson.GetMany(jsonConfigString, "servers", "ak", "password", "consumerId", "topics")
	servers, ak, password := results[0].Array(), results[1].String(), results[2].String()
	consumerId, topics := results[3].String(), results[4].Array()

	s, t := []string{}, []string{}
	s = append(s, servers[0].String())
	t = append(t, topics[0].String())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	client := kafka.New(s, ak, password, true, newLogger)
	consumer, err := client.NewConsumer(consumerId, t, `oldest`)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			fmt.Printf(
				"Topic: %s, Key: %s, Partition: %d, Offset: %d, Content: %s \n", msg.Topic(),
				msg.Key(), msg.Partition(), msg.Offset(), string(msg.Value()),
			)
			consumer.Commit(msg)
		case <-signals:
			fmt.Println("Stop consumer server...")
			consumer.Close()
			return
		}
	}
}
