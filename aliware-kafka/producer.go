package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

type Producer struct {
	topic         string
	producer      sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	logger        *log.Logger
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

	clientCertPool, err := c.AppendValidateCertificate()

	mqConfig.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf(
			"Kafka producer config invalidate. servers: %v. ak: %s, pwd: %s, err: %v",
			c.servers, c.accessKey, c.password, err,
		)
		c.logger.Println(msg)
		panic(msg)
	}
	return mqConfig
}

func (c *Client) NewProducer(topic string) *Producer {
	mgConfig := c.initProducer()
	producer, err := sarama.NewSyncProducer(c.servers, mgConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		c.logger.Println(msg)
		panic(msg)
	}
	return &Producer{topic: topic, producer: producer, logger: c.logger}
}

func (p *Producer) Produce(key string, content string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
	partition, offset, err = p.producer.SendMessage(msg)
	if err != nil {
		msg := fmt.Sprintf("Kafka send message error. topic: %v. key: %v. content: %v", p.topic, key, content)
		p.logger.Println(msg)
		return 0, 0, err
	}
	return partition, offset, nil
}
