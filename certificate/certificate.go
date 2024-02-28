package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

type Certificate struct {
	ca []byte
}

func New(caFile string) (*Certificate, error) {
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	return &Certificate{ca: ca}, nil
}

func (c *Certificate) Server(certFile, keyFile string) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(c.ca); !ok {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}), nil
}

func (c *Certificate) Client(serverName, certFile, keyFile string) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(c.ca); !ok {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:   serverName, //校验证书中的DNS,为空校验IP
		ClientCAs:    certPool,
		Certificates: []tls.Certificate{cert},
	}), nil
}
