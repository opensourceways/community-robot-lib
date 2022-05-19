package mq

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

type MQConfig struct {
	// Addresses of the mq cluster
	Addresses []string `json:"addresses,omitempty"`

	TLSConfig
}

type TLSConfig struct {
	// CertFile the optional certificate file for client authentication
	CertFile string `json:"cert_file,omitempty"`

	// KeyFile the optional key file for client authentication
	KeyFile string `json:"key_file,omitempty"`

	// CAFile  the optional certificate authority file for TLS client authentication
	CAFile string `json:"ca_file,omitempty"`

	// VerifySSL optional verify ssl certificates chain
	VerifySSL bool `json:"verify_ssl,omitempty"`
}

func (tc *TLSConfig) TLSConfig() (t *tls.Config, err error) {
	if tc.CertFile != "" && tc.KeyFile != "" && tc.CAFile != "" {
		cert, err := tls.LoadX509KeyPair(tc.CertFile, tc.KeyFile)
		if err != nil {
			return nil, err
		}

		caCert, err := ioutil.ReadFile(tc.CAFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: tc.VerifySSL,
		}
	}

	return
}
