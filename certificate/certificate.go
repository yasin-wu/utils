package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"os"
	"strings"
)

type Certificate struct {
	caFile string
}

func New(caFile string) *Certificate {
	return &Certificate{caFile: caFile}
}

func (c *Certificate) Server(certFile, keyFile string) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(c.caFile)
	if err != nil {
		return nil, err
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}), nil
}

func (c *Certificate) Client(certFile, keyFile string, serverName ...string) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(c.caFile)
	if err != nil {
		return nil, err
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}

	name := "localhost"
	if len(serverName) > 0 {
		name = strings.Join(serverName, ",")
	}
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   name,
		ClientCAs:    certPool,
	}), nil
}
