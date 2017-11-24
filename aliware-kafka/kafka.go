package kafka

import (
	"crypto/x509"
	"fmt"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"log"
)

const CertFilePath = "./ca-cert"

type Client struct {
	servers   []string
	accessKey string
	password  string
	debug     bool
	logger    *log.Logger
}

func New(servers []string, accessKey, password string, debug bool, logger *log.Logger) *Client {
	fmt.Println("init kafka client")
	if debug {
		sarama.Logger = logger
	}
	return &Client{servers: servers, accessKey: accessKey, password: password, debug: debug, logger: logger}
}

func (c *Client) AppendValidateCertificate() (*x509.CertPool, error) {
	certBytes, err := ioutil.ReadFile(CertFilePath)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka producer failed to parse root certificate")
	}
	return clientCertPool, err
}
