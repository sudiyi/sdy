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
}

type kafkaConsumerMessageWrapper struct {
	message *sarama.ConsumerMessage
}

func (wrapper *kafkaConsumerMessageWrapper) Key() []byte {
	return wrapper.message.Key
}

func (wrapper *kafkaConsumerMessageWrapper) Topic() string {
	return wrapper.message.Topic
}
func (wrapper *kafkaConsumerMessageWrapper) Value() []byte {
	return wrapper.message.Value
}

func (wrapper *kafkaConsumerMessageWrapper) Offset() int64 {
	return wrapper.message.Offset
}

func (wrapper *kafkaConsumerMessageWrapper) Partition() int32 {
	return wrapper.message.Partition
}

func (c *Client) initBasicWithAliwareCertificateConsumer(offset string) *cluster.Config {
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
		log.Fatalln("Offset should be `oldest` or `newest`")
	}
	return initialOffset
}

func (c *Client) initConsumer(offset string) (*cluster.Config, error) {
	switch c.encrypt {
	case "aliware":
		return c.initWithAliwareCertificateConsumer(offset)
	default:
		return c.initWithNonCertificateConsumer(offset)
	}
}

func (c *Client) initWithAliwareCertificateConsumer(offset string) (*cluster.Config, error) {
	clusterCfg := c.initBasicWithAliwareCertificateConsumer(offset)
	clientCertPool, err := c.AppendValidateCertificate()
	if err != nil {
		return nil, err
	}
	clusterCfg.Net.TLS.Config = &tls.Config{
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}
	if err := clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		log.Println(msg)
		panic(msg)
	}
	return clusterCfg, nil
}

func (c *Client) initWithNonCertificateConsumer(offset string) (*cluster.Config, error) {
	clusterCfg := cluster.NewConfig()

	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = c.getInitialOffset(offset)
	clusterCfg.Group.Return.Notifications = true
	clusterCfg.Version = sarama.V0_10_0_0
	return clusterCfg, nil
}

func (c *Client) NewConsumer(consumerId string, topics []string, offset string) (*AliyunConsumer, error) {
	clusterCfg, err := c.initConsumer(offset)
	if err != nil {
		return nil, err
	}
	consumer, err := cluster.NewConsumer(c.servers, consumerId, topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		log.Println(msg)
		return nil, err
	}
	aliYun := &AliyunConsumer{
		consumer: consumer,
		messages: make(chan Message),
	}

	go aliYun.run()

	return aliYun, nil
}

func (consumer *AliyunConsumer) run() {
	for {
		select {
		case msg, more := <-consumer.consumer.Messages():
			if more {
				consumer.messages <- &kafkaConsumerMessageWrapper{msg}
			} else {
				// maybe get the error about "close of closed channel"
				//close(consumer.messages)
			}
		case notify, more := <-consumer.consumer.Notifications():
			if more {
				log.Println("Kafka consumer rebalanced: %v", notify)
			}
		case err, more := <-consumer.consumer.Errors():
			if more {
				log.Printf("Errors: %s\n", err.Error())
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
	consumer.consumer.MarkOffset(message.(*kafkaConsumerMessageWrapper).message, "")
}
