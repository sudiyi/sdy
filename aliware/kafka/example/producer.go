package main

import (
	"fmt"
	"github.com/sudiyi/sdy/aliware/kafka"
	"github.com/tidwall/gjson"
)

var jsonString = `
{
  "topics": ["your topic"],
  "servers": ["kafka-ons-internet.aliyun.com:8080"],
  "ak": "Access Key",
  "password": "password",
  "consumerId": "your consumer id",
}
`

func main() {
	results := gjson.GetMany(jsonString, "servers", "ak", "password", "consumerId", "topics")
	servers, ak, password := results[0].Array(), results[1].String(), results[2].String()
	s := []string{}
	s = append(s, servers[0].String())
	_, topics := results[3].String(), results[4].Array()

	// litte concurrency
	syncProducer(s, ak, password, topics[0].String())

	// high concurrency
	asyncProducer(s, ak, password, topics[0].String())
}

func asyncProducer(s []string, ak, password string, topic string) {
	client := kafka.New(s, ak, password, false)
	p, _ := client.NewAsyncProducer(topic)
	defer p.AsyncClose()

	var i int

	for {
		for _, key := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"} {
			//for _, key := range []string{"A"} {
			p.AsyncProduce(key+"-GOOD", `{"a": 1, "b": [{"a": 1}]}`)
		}
		i += 1
		if i > 0 {
			break
		}
	}

}

func syncProducer(s []string, ak, password string, topic string) {
	client := kafka.New(s, ak, password, false)
	p, _ := client.NewProducer(topic)

	var i int

	for {
		for _, key := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"} {
			//for _, key := range []string{"A"} {
			fmt.Println(p.Produce(key+"-NEW-TEST", `{"a": 1, "b": [{"a": 1}]}`))
		}
		i += 1
		if i > 0 {
			break
		}
	}
}
