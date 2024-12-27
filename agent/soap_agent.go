package agent

import (
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

// type SoapAgent struct {
// 	Config     *Config
// 	SopaSender *SoapSender
// }

// func NewSoapAgent() *SoapAgent {
// 	soapAgent := SoapAgent{}
// 	return &soapAgent
// }

func (soapSender *SoapSender) ReserveSender(siteKey string, filename string) error {
	req, err := soapSender.MakeSoapRequestMsg(
		"reserveSender",
		map[string]string{
			"siteKey":  siteKey,
			"filename": filename,
			"size":     "0",
		},
	)
	if err != nil {
		errors.Wrap(err, "reserve sender request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		errors.Wrap(err, "make request to reserve sender")
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		return errors.Wrap(err, "soap call to reserve sender")
	}
	log.Infof("response:%v\n", response)
	return nil
}

func (soapSender *SoapSender) SendData(siteKey, filename, zipPath string) (string, error) {

	req, err := soapSender.MakeSoapRequestMsg(
		"sendData",
		map[string]string{
			"siteKey":  siteKey,
			"filename": filename,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "send data request generation")
	}
	httpReq, err := soapSender.MakeSoapRequestWithAttachment(req, zipPath)
	if err != nil {
		return "", errors.Wrap(err, "make request to send data")
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		// log.Fatalln(err)
		return "", errors.Wrap(err, "soap call to send data")
	}
	log.Infof("send data response '%v'\n", response)
	return response, nil
}

func (soapSender *SoapSender) SendMessage(siteKey, hostname, severity, message string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"sendMessage",
		map[string]string{
			"siteKey":  siteKey,
			"hostname": hostname,
			"severity": severity,
			"message":  message,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "send message request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		errors.Wrap(err, "make request to send message")
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		// log.Fatalln(err)
		return "", errors.Wrap(err, "soap call to send message")
	}
	log.Infof("send message response '%v'\n", response)
	return response, nil
}

func (soapSender *SoapSender) DownloadCertificate(siteKey, hostname, timestamp, netConfigPath string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"downloadCertificate",
		map[string]string{
			"siteKey":   siteKey,
			"hostname":  hostname,
			"timestamp": timestamp,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "download certificate request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		errors.Wrap(err, "make request to download certificate")
	}
	soapSender.WithAttachedFilePath(netConfigPath)
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		// log.Fatalln(err)
		return "", errors.Wrap(err, "soap call to download certificate")
	}
	log.Infof("download certificate response '%v'\n", response)
	return response, nil
}

func (soapSender *SoapSender) ReserveFileSender(onOff string, waitSec int) error {

	return nil
}

func (soapSender *SoapSender) SendZipData(filename string) error {
	return nil
}

func (soapSender *SoapSender) DownloadConfigFilePM(filename string) error {
	return nil
}

func (soapSender *SoapSender) GetLatestBuild(moduleTag string, majorVer string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"getLatestBuild",
		map[string]string{
			"moduleTag": moduleTag,
			"majorVer":  majorVer,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "reserve sender request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		return "", errors.Wrap(err, "make request to reserve sender")
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "soap call to reserve sender")
	}
	log.Infof("get latest build response '%s'", response)

	return response, nil
}

func (soapSender *SoapSender) RegistAgent(siteKey, hostname, accessKey, netConfigPath string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"registAgent",
		map[string]string{
			"siteKey":   siteKey,
			"hostname":  hostname,
			"accessKey": accessKey,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "regist agent request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		return "", errors.Wrap(err, "make request to regist agent")
	}
	soapSender.WithAttachedFilePath(netConfigPath)
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "soap call to regist agent")
	}
	log.Infof("regist agent response '%s'", response)

	return response, nil
}

func (soapSender *SoapSender) DownloadUpdateModule(moduleTag, majorVer, build, downloadPath string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"downloadUpdateModule",
		map[string]string{
			"moduleTag": moduleTag,
			"majorVer":  majorVer,
			"build":     build,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "download update module request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		return "", errors.Wrap(err, "make request to douwnload update module")
	}
	soapSender.WithAttachedFilePath(downloadPath)
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "soap call to download update module")
	}
	log.Infof("download update module response '%s'", response)

	return response, nil
}

// TODO 引数に GPFSetupConfig *setup を追加
func (soapSender *SoapSender) CheckHostStatus(siteKey, hostname, accessKey string) (string, error) {
	req, err := soapSender.MakeSoapRequestMsg(
		"checkAgent",
		map[string]string{
			"siteKey":   siteKey,
			"hostname":  hostname,
			"accessKey": accessKey,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "check host status request generation")
	}
	httpReq, err := soapSender.MakeSoapRequest(req)
	if err != nil {
		return "", errors.Wrap(err, "make request to check host status")
	}
	response, err := soapSender.soapCall(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "soap call to check host status")
	}
	log.Infof("check host status response '%s'", response)

	return response, nil
}
