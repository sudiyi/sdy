package kafka

import (
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

const AliyunCertificate = `
-----BEGIN CERTIFICATE-----
MIIDPDCCAqWgAwIBAgIJAMRsb0DLM1fsMA0GCSqGSIb3DQEBBQUAMHIxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJIWjELMAkGA1UEBxMCSFoxCzAJBgNVBAoTAkFCMRAw
DgYDVQQDEwdLYWZrYUNBMSowKAYJKoZIhvcNAQkBFht6aGVuZG9uZ2xpdS5semRA
YWxpYmFiYS5jb20wIBcNMTcwMzA5MTI1MDUyWhgPMjEwMTAyMTcxMjUwNTJaMHIx
CzAJBgNVBAYTAkNOMQswCQYDVQQIEwJIWjELMAkGA1UEBxMCSFoxCzAJBgNVBAoT
AkFCMRAwDgYDVQQDEwdLYWZrYUNBMSowKAYJKoZIhvcNAQkBFht6aGVuZG9uZ2xp
dS5semRAYWxpYmFiYS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALZV
bbIO1ULQQN853BTBgRfPiRJaAOWf38u8GC0TNp/E9qtI88A+79ywAP17k5WYJ7XS
wXMOJ3h1qkQT2TYJVetZ6E69CUJq4BsOvNlNRvmnW6eFymh5QZsEz2MTooxJjVjC
JQPlI2XRDjIrTVYEQWUDxj2JhB8VVPEed+6u4KQVAgMBAAGjgdcwgdQwHQYDVR0O
BBYEFHFlOoiqQxXanVi2GUoDiKDD33ujMIGkBgNVHSMEgZwwgZmAFHFlOoiqQxXa
nVi2GUoDiKDD33ujoXakdDByMQswCQYDVQQGEwJDTjELMAkGA1UECBMCSFoxCzAJ
BgNVBAcTAkhaMQswCQYDVQQKEwJBQjEQMA4GA1UEAxMHS2Fma2FDQTEqMCgGCSqG
SIb3DQEJARYbemhlbmRvbmdsaXUubHpkQGFsaWJhYmEuY29tggkAxGxvQMszV+ww
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOBgQBTSz04p0AJXKl30sHw+UM/
/k1jGFJzI5p0Z6l2JzKQYPP3PfE/biE8/rmiGYEenNqWNy1ZSniEHwa8L/Ux98ci
4H0ZSpUrMo2+6bfuNW9X35CFPp5vYYJqftilJBKIJX3C3J1ruOuBR28UxE42xx4K
pQ70wChNi914c4B+SxkGUg==
-----END CERTIFICATE-----
`

type Client struct {
	servers   []string
	accessKey string
	password  string
	debug     bool
	encrypt   string
}

func New(servers []string, debug bool) *Client {
	fmt.Println("init kafka client")
	if debug {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}
	return &Client{servers: servers}
}

func (c *Client) EncryptByAliware(accessKey, password string) *Client {
	c.accessKey = accessKey
	c.password = password
	c.encrypt = "aliware"
	return c
}

func (c *Client) AppendValidateCertificate() (*x509.CertPool, error) {
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM([]byte(AliyunCertificate))
	if !ok {
		return nil, errors.New("kafka producer failed to parse root certificate")
	}
	return clientCertPool, nil
}
