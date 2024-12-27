package agent

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (soapSender *SoapSender) MakeTransportWithServerAuthSSL(c *Config) (*http.Transport, error) {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "init ca file")
	}
	caCert, err := ioutil.ReadFile(c.CacertFile)
	if err != nil {
		return nil, errors.Wrap(err, "read cafile to make soap ssl transport")
	}
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.Wrap(err, "failed to add ca cert")
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.Wrap(err, "invalid default transport")
	}

	transport := defaultTransport.Clone()

	transport.TLSClientConfig = &tls.Config{
		RootCAs:    caCertPool,
		ServerName: soapSender.ServerIP,
	}
	soapSender.Transport = transport
	return transport, nil
}

func (soapSender *SoapSender) MakeTransportWithClientAuthSSL(c *Config) (*http.Transport, error) {
	caCertFile := c.CacertFile
	clientCertFile := c.ClcertFile
	clientKeyFile := c.ClkeyFile

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"Error creating x509 keypair from client cert file %s and client key file %s",
			clientCertFile,
			clientKeyFile,
		)
	}
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.Wrap(err, "invalid default transport")
	}
	transport := defaultTransport.Clone()

	transport.TLSClientConfig = &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		ServerName:   soapSender.ServerIP,
	}
	soapSender.Transport = transport
	return transport, nil
}
