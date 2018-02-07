package main

import (
	"fmt"
	"github.com/sudiyi/sdy/aliware/kafka"
	"github.com/tidwall/gjson"
)

var demoJsonString = `
{
  "topics": ["demo"],
  "servers": ["kafka1:9092"],
  "consumerId": "demo-consumer-group",
}
`

func main() {
	results := gjson.GetMany(demoJsonString, "servers", "consumerId", "topics")
	servers := results[0].Array()
	var s []string
	s = append(s, servers[0].String())
	_, topics := results[1].String(), results[2].Array()

	//// litte concurrency
	syncNonProducer(s, topics[0].String())
	//
	//// high concurrency
	asyncNonProducer(s, topics[0].String())
}

func syncNonProducer(s []string, topic string) {
	client := kafka.New(s, true)
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

func asyncNonProducer(s []string, topic string) {
	client := kafka.New(s, true)
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
