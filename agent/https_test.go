package agent

import (
	"net/url"
	"testing"
)

func TestGetTSLConfig(t *testing.T) {
	config := NewConfig("../testdata/ptune", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	if config.SvcertFile != "../testdata/ptune/network/server/server.crt" ||
		config.SvkeyFile != "../testdata/ptune/network/server/server.key" ||
		config.SvcacertFile != "../testdata/ptune/network/server/ca.crt" {
		t.Error("tls config")
	}
	if config.Schedule.WebServiceEnable != true ||
		config.Schedule.WebServiceUrl != "https://0.0.0.0:59443" {
		t.Error("tls config2")
	}
	u, err := url.Parse(config.Schedule.WebServiceUrl)
	if err != nil {
		t.Error(err)
	}
	t.Log(u.Scheme)
}

func TestMakeTLSConfigClientAuth(t *testing.T) {
	config := NewConfig("../testdata/ptune", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	tlsConfig, err := MakeTLSConfigClientAuth(config)
	if err != nil {
		t.Error(err)
	}
	t.Log(tlsConfig)
}

func TestMakeTLSConfigServerAuth(t *testing.T) {
	t.Skip("Server TLS Config")
	config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	tlsConfig, err := MakeTLSConfigServerAuth(config)
	if err != nil {
		t.Error(err)
	}
	t.Log(tlsConfig)
}
