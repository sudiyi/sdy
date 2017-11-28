package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

const (
	DefaultPartition int32 = 0
	DefaultOffset    int64 = 0
)

type Producer struct {
	topic         string
	producer      sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
}

func (c *Client) initConfigForProducer() *sarama.Config {
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = c.accessKey
	mqConfig.Net.SASL.Password = c.password
	mqConfig.Net.SASL.Handshake = true

	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Errors = true
	mqConfig.Producer.Return.Successes = true
	return mqConfig
}

func (c *Client) initProducer() *sarama.Config {
	mqConfig := c.initConfigForProducer()

	clientCertPool := c.AppendValidateCertificate()

	mqConfig.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	if err := mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf(
			"Kafka producer config invalidate. servers: %v. ak: %s, pwd: %s, err: %v",
			c.servers, c.accessKey, c.password, err,
		)
		log.Println(msg)
		panic(msg)
	}
	return mqConfig
}

func (c *Client) NewProducer(topic string) *Producer {
	mgConfig := c.initProducer()
	producer, err := sarama.NewSyncProducer(c.servers, mgConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		log.Println(msg)
		panic(msg)
	}
	return &Producer{topic: topic, producer: producer}
}

// sync producer, use for little concurrency
func (p *Producer) Produce(key string, content string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
	partition, offset, err = p.producer.SendMessage(msg)
	if err != nil {
		msg := fmt.Sprintf("Kafka send message error. topic: %v. key: %v. content: %v", p.topic, key, content)
		log.Println(msg)
		return DefaultPartition, DefaultOffset, err
	}
	return partition, offset, nil
}

// async producer, use for many concurrency
func (c *Client) NewAsyncProducer(topic string) *Producer {
	mgConfig := c.initProducer()
	asyncProducer, err := sarama.NewAsyncProducer(c.servers, mgConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak async producer create fail. err: %v", err)
		log.Println(msg)
		panic(msg)
	}
	return &Producer{topic: topic, AsyncProducer: asyncProducer}
}

func (p *Producer) AsyncProduce(key string, content string) {
	go func() {
		errors := p.AsyncProducer.Errors()
		success := p.AsyncProducer.Successes()

		for {
			select {
			case err, ok := <-errors:
				if ok {
					log.Fatalln("FAILURE:", err)
				}
			case message, ok := <-success:
				if ok {
					log.Printf(
						"Topic: %s, Key: %s, Partition: %d, Offset: %d \n", message.Topic,
						message.Key, message.Partition, message.Offset,
					)
				}
			}
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
	p.AsyncProducer.Input() <- msg
}

func (p *Producer) AsyncClose() {
	p.AsyncProducer.Close()
}
