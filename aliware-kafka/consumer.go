package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

type Consumer struct {
	consumer   *cluster.Consumer
	consumerId string
	topics     []string
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

func (c *Client) NewConsumer(consumerId string, topics []string, offset string) *Consumer {
	clusterCfg := c.initConsumer(offset)
	consumer, err := cluster.NewConsumer(c.servers, consumerId, topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		c.logger.Println(msg)
		panic(msg)
	}
	return &Consumer{consumer: consumer, consumerId: consumerId, topics: topics}
}
