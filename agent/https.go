// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package https allows the implementation of TLS.
package agent

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const ClientAuth = "RequireAndVerifyClientCert"

func MakeTLSConfigClientAuth(c *Config) (*tls.Config, error) {
	return MakeTLSConfig(c, "RequireAndVerifyClientCert")
	// cfg := &tls.Config{}
	// if len(c.SvcertFile) == 0 {
	// 	return nil, errors.New("missing SvcertFile")
	// }
	// if len(c.SvkeyFile) == 0 {
	// 	return nil, errors.New("missing SvkeyFile")
	// }
	// loadCert := func() (*tls.Certificate, error) {
	// 	cert, err := tls.LoadX509KeyPair(c.SvcertFile, c.SvkeyFile)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "failed to load X509KeyPair")
	// 	}
	// 	return &cert, nil
	// }
	// // Confirm that certificate and key paths are valid.
	// if _, err := loadCert(); err != nil {
	// 	return nil, err
	// }
	// cfg.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	// 	return loadCert()
	// }

	// if len(c.SvcacertFile) > 0 {
	// 	clientCAPool := x509.NewCertPool()
	// 	clientCAFile, err := ioutil.ReadFile(c.SvcacertFile)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	clientCAPool.AppendCertsFromPEM(clientCAFile)
	// 	cfg.ClientCAs = clientCAPool
	// }
	// if len(ClientAuth) > 0 {
	// 	switch s := (ClientAuth); s {
	// 	case "NoClientCert":
	// 		cfg.ClientAuth = tls.NoClientCert
	// 	case "RequestClientCert":
	// 		cfg.ClientAuth = tls.RequestClientCert
	// 	case "RequireClientCert":
	// 		cfg.ClientAuth = tls.RequireAnyClientCert
	// 	case "VerifyClientCertIfGiven":
	// 		cfg.ClientAuth = tls.VerifyClientCertIfGiven
	// 	case "RequireAndVerifyClientCert":
	// 		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	// 	case "":
	// 		cfg.ClientAuth = tls.NoClientCert
	// 	default:
	// 		return nil, errors.New("Invalid ClientAuth: " + s)
	// 	}
	// }
	// if len(c.SvcacertFile) > 0 && cfg.ClientAuth == tls.NoClientCert {
	// 	return nil, errors.New("Client CA's have been configured without a Client Auth Policy")
	// }
	// return cfg, nil
}

func MakeTLSConfigServerAuth(c *Config) (*tls.Config, error) {
	return MakeTLSConfig(c, "")
}

func MakeTLSConfig(c *Config, clientAuth string) (*tls.Config, error) {
	cfg := &tls.Config{}
	if len(c.SvcertFile) == 0 {
		return nil, errors.New("missing SvcertFile")
	}
	if len(c.SvkeyFile) == 0 {
		return nil, errors.New("missing SvkeyFile")
	}
	loadCert := func() (*tls.Certificate, error) {
		cert, err := tls.LoadX509KeyPair(c.SvcertFile, c.SvkeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load X509KeyPair")
		}
		return &cert, nil
	}
	// Confirm that certificate and key paths are valid.
	if _, err := loadCert(); err != nil {
		return nil, err
	}
	cfg.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return loadCert()
	}

	if len(c.SvcacertFile) > 0 {
		clientCAPool := x509.NewCertPool()
		clientCAFile, err := ioutil.ReadFile(c.SvcacertFile)
		if err != nil {
			return nil, err
		}
		clientCAPool.AppendCertsFromPEM(clientCAFile)
		cfg.ClientCAs = clientCAPool
	}
	if len(clientAuth) > 0 {
		switch s := (clientAuth); s {
		case "NoClientCert":
			cfg.ClientAuth = tls.NoClientCert
		case "RequestClientCert":
			cfg.ClientAuth = tls.RequestClientCert
		case "RequireClientCert":
			cfg.ClientAuth = tls.RequireAnyClientCert
		case "VerifyClientCertIfGiven":
			cfg.ClientAuth = tls.VerifyClientCertIfGiven
		case "RequireAndVerifyClientCert":
			cfg.ClientAuth = tls.RequireAndVerifyClientCert
		case "":
			cfg.ClientAuth = tls.NoClientCert
		default:
			return nil, errors.New("Invalid ClientAuth: " + s)
		}
	}
	if len(c.SvcacertFile) > 0 && cfg.ClientAuth == tls.NoClientCert {
		return nil, errors.New("Client CA's have been configured without a Client Auth Policy")
	}
	return cfg, nil
}

// Listen starts the server on the given address. If tlsConfigPath isn't empty the server connection will be started using TLS.
func Listen(server *http.Server, config *Config) error {
	u, err := url.Parse(config.Schedule.WebServiceUrl)
	if err != nil {
		return errors.Wrap(err, "prepare listen port")
	}
	if u.Scheme == "http" {
		return server.ListenAndServe()
	}
	server.TLSConfig, err = MakeTLSConfigClientAuth(config)
	if err != nil {
		return err
	}
	// Set the GetConfigForClient method of the HTTPS server so that the config
	// and certs are reloaded on new connections.
	server.TLSConfig.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		return MakeTLSConfigClientAuth(config)
	}
	return server.ListenAndServeTLS("", "")
}
