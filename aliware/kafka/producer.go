package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
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

func (c *Client) initProducer() (*sarama.Config, error) {
	switch c.encrypt {
	case "aliware":
		return c.initWithAliwareCertificateProducer()
	default:
		return c.initWithNonCertificateProducer()
	}
}

func (c *Client) initBasicWithAliwareCertificateProducer() *sarama.Config {
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = c.accessKey
	mqConfig.Net.SASL.Password = c.password
	mqConfig.Net.SASL.Handshake = true
	mqConfig.Net.ReadTimeout = 10 * time.Second
	mqConfig.Net.DialTimeout = 10 * time.Second
	mqConfig.Net.WriteTimeout = 10 * time.Second

	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Errors = true
	mqConfig.Producer.Return.Successes = true
	mqConfig.Producer.Retry.Backoff = 10 * time.Second
	mqConfig.Producer.Retry.Max = 3
	
	mqConfig.Metadata.Retry.Max = 1
	mqConfig.Metadata.Retry.Backoff = 10 * time.Second
	mqConfig.Metadata.RefreshFrequency = 15 * time.Minute
	return mqConfig
}

func (c *Client) initWithAliwareCertificateProducer() (*sarama.Config, error) {
	mqConfig := c.initBasicWithAliwareCertificateProducer()
	clientCertPool, err := c.AppendValidateCertificate()
	if err != nil {
		return nil, err
	}
	mqConfig.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	err = mqConfig.Validate()
	if err != nil {
		msg := fmt.Sprintf(
			"Kafka producer config invalidate. servers: %v. ak: %s, pwd: %s, err: %v",
			c.servers, c.accessKey, c.password, err,
		)
		log.Println(msg)
	}
	return mqConfig, err
}

func (c *Client) initWithNonCertificateProducer() (*sarama.Config, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner
	return config, nil
}

func (c *Client) NewProducer(topic string) (*Producer, error) {
	mgConfig, err := c.initProducer()
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer(c.servers, mgConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		log.Println(msg)
	}
	return &Producer{topic: topic, producer: producer}, err
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
func (c *Client) NewAsyncProducer(topic string) (*Producer, error) {
	mgConfig, err := c.initProducer()
	if err != nil {
		return nil, err
	}
	asyncProducer, err := sarama.NewAsyncProducer(c.servers, mgConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak async producer create fail. err: %v", err)
		log.Println(msg)
	}
	return &Producer{topic: topic, AsyncProducer: asyncProducer}, err
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
