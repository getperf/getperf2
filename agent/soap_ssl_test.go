package agent

import (
	"testing"
)

var TEST_WEB_SERVICE_HOST = "192.168.0.65" // "10.45.50.210" // "192.168.133.128"

func TestMakeTransportWithServerAuthSSL(t *testing.T) {
	config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	soapSender, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 57443)
	if err != nil {
		t.Error(err)
	}
	_, err := soapSender.MakeTransportWithServerAuthSSL(config)
	if err != nil {
		t.Error(err)
	}
	t.Logf("soap ssl transport : %v", soapSender.Transport)
}

func TestMakeTransportWithClientAuthSSL(t *testing.T) {
	config := NewConfig("../testdata/ptune-3.0", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	soapSender, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 57443)
	if err != nil {
		t.Error(err)
	}
	_, err := soapSender.MakeTransportWithClientAuthSSL(config)
	if err != nil {
		t.Error(err)
	}
	t.Logf("soap ssl transport : %v", soapSender.Transport)
}
