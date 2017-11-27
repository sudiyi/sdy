package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
)

// AliyunConsumer .
type AliyunConsumer struct {
	consumer *cluster.Consumer
	messages chan Message
	logger   *log.Logger
}

type kafkaConsumerMessageWraper struct {
	message *sarama.ConsumerMessage
}

func (wraper *kafkaConsumerMessageWraper) Key() []byte {
	return wraper.message.Key
}

func (wraper *kafkaConsumerMessageWraper) Topic() string {
	return wraper.message.Topic
}
func (wraper *kafkaConsumerMessageWraper) Value() []byte {
	return wraper.message.Value
}

func (wraper *kafkaConsumerMessageWraper) Offset() int64 {
	return wraper.message.Offset
}

func (wraper *kafkaConsumerMessageWraper) Partition() int32 {
	return wraper.message.Partition
}

func (c *Client) initConfigForConsumer(offset string) *cluster.Config {
	clusterCfg := cluster.NewConfig()

	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = c.accessKey
	clusterCfg.Net.SASL.Password = c.password
	clusterCfg.Net.SASL.Handshake = true

	clusterCfg.Net.TLS.Enable = true

	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = c.getInitialOffset(offset)
	clusterCfg.Group.Return.Notifications = true
	clusterCfg.Version = sarama.V0_10_0_0
	return clusterCfg
}

func (c *Client) getInitialOffset(offset string) int64 {
	var initialOffset int64
	switch offset {
	case "oldest":
		initialOffset = sarama.OffsetOldest
	case "newest":
		initialOffset = sarama.OffsetNewest
	default:
		c.logger.Fatalln("Offset should be `oldest` or `newest`")
	}
	return initialOffset
}

func (c *Client) initConsumer(offset string) *cluster.Config {
	clusterCfg := c.initConfigForConsumer(offset)
	clientCertPool, err := c.AppendValidateCertificate()
	clusterCfg.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		c.logger.Println(msg)
		panic(msg)
	}
	return clusterCfg
}

func (c *Client) NewConsumer(consumerId string, topics []string, offset string) (*AliyunConsumer, error) {
	clusterCfg := c.initConsumer(offset)
	consumer, err := cluster.NewConsumer(c.servers, consumerId, topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		c.logger.Println(msg)
		return nil, err
	}
	aliyun := &AliyunConsumer{
		consumer: consumer,
		messages: make(chan Message),
		logger:   c.logger,
	}

	go aliyun.run()

	return aliyun, nil
}

func (consumer *AliyunConsumer) run() {
	for {
		select {
		case msg, more := <-consumer.consumer.Messages():
			if more {
				consumer.messages <- &kafkaConsumerMessageWraper{msg}
			} else {
				close(consumer.messages)
			}
		case notify, more := <-consumer.consumer.Notifications():
			if more {
				consumer.logger.Println("Kafka consumer rebalanced: %v", notify)
			}
		}
	}
}

// Close .
func (consumer *AliyunConsumer) Close() {
	consumer.consumer.Close()
}

// Messages return message chan
func (consumer *AliyunConsumer) Messages() <-chan Message {
	return consumer.messages
}

// Errors return error chan
func (consumer *AliyunConsumer) Errors() <-chan error {
	return consumer.consumer.Errors()
}

// Commit commit current handle message as consumed
func (consumer *AliyunConsumer) Commit(message Message) {
	consumer.consumer.MarkOffset(message.(*kafkaConsumerMessageWraper).message, "")
}
