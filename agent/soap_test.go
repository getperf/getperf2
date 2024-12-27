package agent

import (
	// "fmt"

	"runtime"
	"testing"
)

var err error
var soapSender, soapSenderData *SoapSender

func init() {
	// config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
	// config.InitAgent()
	// config.ParseConfigFile(config.SvParamFile)
	// soapSender, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 57443)
	// if err != nil {
	// 	panic(err)
	// }
	// // _, err := soapSender.MakeTransportWithSSL()
	// soapSenderData, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 58443)
	// if err != nil {
	// 	panic(err)
	// }

	if runtime.GOOS != "windows" {
		initSoapSenderAdmin()
		initSoapSenderData()
	}
}

func initSoapSenderAdmin() {
	config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	soapSender, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 57443)
	if err != nil {
		panic(err)
	}
	_, err := soapSender.MakeTransportWithServerAuthSSL(config)
	if err != nil {
		panic(err)
	}
}

func initSoapSenderData() {
	config := NewConfig("../testdata/ptune-3.0", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile(config.SvParamFile)
	soapSenderData, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 58443)
	if err != nil {
		panic(err)
	}
	_, err := soapSenderData.MakeTransportWithClientAuthSSL(config)
	if err != nil {
		panic(err)
	}
}

func TestNewSoapSender(t *testing.T) {
	soapSender, err = NewSoapSender(TEST_WEB_SERVICE_HOST, 57443)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	soapSender.WithTimeout(15).WithAttachedFilePath("/tmp/test.dat")
	t.Logf("soapSender new : %v", soapSender)
}

func TestMakeSoapRequestMsg(t *testing.T) {
	soapRequest := map[string]string{
		"moduleTag": "Windows",
		"majorVer":  "2",
	}
	message, err := soapSender.MakeSoapRequestMsg("getLatestBuild", soapRequest)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	t.Logf("%v", message)
}

func TestMakeSoapSendMessageRequestMsg(t *testing.T) {
	soapRequest := map[string]string{
		"siteKey":  "site1",
		"hostname": "host1",
		"severity": "1",
		"message":  "this is a test",
	}
	msg, err := soapSender.MakeSoapRequestMsg("sendMessage", soapRequest)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	t.Logf("%v", msg)
}

func TestSoapGetLatestBuildRequest(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"getLatestBuild",
		map[string]string{
			"moduleTag": "Linux",
			"majorVer":  "2",
		},
	)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	t.Logf("req:%v", req)
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	t.Logf("request %v", httpReq)
}

func TestSoapCallGetLatestBuild(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"getLatestBuild",
		map[string]string{
			"moduleTag": "Linux",
			"majorVer":  "2",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallSendMessage(t *testing.T) {
	req, err := soapSenderData.MakeSoapRequestMsg(
		"sendMessage",
		map[string]string{
			"siteKey":  "site1",
			"hostname": "host1",
			"severity": "1",
			"message":  "this is a test3",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSenderData.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSenderData.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallSendData(t *testing.T) {
	req, err := soapSenderData.MakeSoapRequestMsg(
		"sendData",
		map[string]string{
			"siteKey":  "site1",
			"filename": "arc_host1__Linux_20230506_0800.zip",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSenderData.MakeSoapRequestWithAttachment(req, "../arc_host1__Linux_20230506_0800.zip")
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSenderData.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallRegistAgent(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"registAgent",
		map[string]string{
			"siteKey":   "site1",
			"hostname":  "host1",
			"accessKey": "81e6011f1c0660a8062dbe4ade4e910d841d36c4",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	soapSender.WithAttachedFilePath("/tmp/test.zip")
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapGetResponseReturn(t *testing.T) {
	msg, err := soapSender.getResponseReturn(`
		<?xml version='1.0' encoding='UTF-8'?><soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"><soapenv:Header xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Action>urn:getLatestBuildResponse</wsa:Action><wsa:RelatesTo>urn:uuid:bfd15f5f-c832-459a-935b-201494336637</wsa:RelatesTo></soapenv:Header><soapenv:Body><ns:getLatestBuildResponse xmlns:ns=\"http://perf.getperf.com\"><ns:return>0</ns:return></ns:getLatestBuildResponse></soapenv:Body></soapenv:Envelope>
		`)
	if err != nil {
		t.Logf("response parse error %s", err)
	}
	t.Logf("response %v", msg)

}

func TestSoapGetResponseReturnError(t *testing.T) {
	msg, err := soapSender.getResponseReturn(`
		Hoge`)
	if err == nil {
		t.Logf("response parse error %s", err)
	}
	t.Logf("errormsg %v", err)
	t.Logf("response %v", msg)

}

func TestSoapCallDownloadCertificate(t *testing.T) {
	req, err := soapSenderData.MakeSoapRequestMsg(
		"downloadCertificate",
		map[string]string{
			"siteKey":   "site1",
			"hostname":  "host1",
			"timestamp": "0",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSenderData.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	soapSenderData.WithAttachedFilePath("/tmp/test.zip")
	response, err := soapSenderData.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallCehckAgent(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"checkAgent",
		map[string]string{
			"siteKey":   "site1",
			"hostname":  "host1",
			"accessKey": "81e6011f1c0660a8062dbe4ade4e910d841d36c4",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallGetLatestVersion(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"getLatestVersion",
		map[string]string{},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallReserveSender(t *testing.T) {
	req, err := soapSenderData.MakeSoapRequestMsg(
		"reserveSender",
		map[string]string{
			"siteKey":  "site1",
			"filename": "arc_host1__Linux_20230506_0800.zip",
			"size":     "0",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSenderData.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	response, err := soapSenderData.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}

func TestSoapCallDownloadUpdateModule(t *testing.T) {
	req, err := soapSender.MakeSoapRequestMsg(
		"downloadUpdateModule",
		map[string]string{
			"moduleTag": "Linux",
			"majorVer":  "2",
			"build":     "1",
		},
	)
	if err != nil {
		t.Errorf("Some problem occurred in request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		t.Errorf("check request %s", err)
	}
	soapSender.WithAttachedFilePath("/tmp/test.zip")
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		t.Logf("soap call error %s", err)
	}
	t.Logf("response %v", response)
}
