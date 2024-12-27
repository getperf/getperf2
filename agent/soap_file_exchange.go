package agent

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (config *Config) RunFileTransferDownloadLicense() error {
	return config.RunFileTransfer(false, "")
}

func (config *Config) RunFileTransfer(sendEnable bool, zipFile string) error {
	if err := config.SetupHttpProxy(); err != nil {
		log.Fatal("set http proxy ", err)
	}

	schedule := config.Schedule
	u, err := url.ParseRequestURI(schedule.UrlPM)
	if err != nil {
		return errors.Wrap(err, "parse URL_PM")
	}
	portNum, _ := strconv.Atoi(u.Port())
	log.Debugf("set web service host: %v, port: %v", u.Hostname(), u.Port())
	soapSender, err := NewSoapSender(u.Hostname(), portNum)
	if err != nil {
		return errors.Wrap(err, "initialize soap sender")
	}
	_, err = soapSender.MakeTransportWithClientAuthSSL(config)
	if err != nil {
		return errors.Wrap(err, "set client cert config")
	}
	zipPath := config.GetArchivefilePath(zipFile)
	log.Debugf("set zip path %s", zipPath)
	if sendEnable {
		for retry := 0; retry < DEFAULT_SOAP_RETRY; retry++ {
			log.Infof("send %v [%d/%d]", zipPath, retry+1, DEFAULT_SOAP_RETRY)
			result, err := soapSender.SendData(schedule.SiteKey, zipFile, zipPath)
			if err != nil {
				log.Errorf("file send %s", err)
				if retry < DEFAULT_SOAP_RETRY {
					time.Sleep(time.Duration(DEFAULT_SOAP_TIMEOUT) * time.Second)
				}
				continue
			} else {
				log.Infof("file send end %s", result)
			}
			break
		}
	} else {
		sslPath := filepath.Join(config.WorkCommonDir, "sslconf.zip")
		for retry := 0; retry < DEFAULT_SOAP_RETRY; retry++ {
			log.Infof("get %v [%d/%d]", sslPath, retry+1, DEFAULT_SOAP_RETRY)
			result, err := soapSender.DownloadCertificate(schedule.SiteKey, config.Host, "0", sslPath)
			if err != nil {
				log.Errorf("file get %s", err)
				if retry < DEFAULT_SOAP_RETRY {
					time.Sleep(time.Duration(DEFAULT_SOAP_TIMEOUT) * time.Second)
				}
				continue
			} else {
				log.Infof("file get end %s", result)
			}
			break
		}
	}
	return err
}

func RunSender(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
	var usage = `getperfsoap [--send(-s)|--get(-g)] [--config(-c) getperf.cfg] 
	        filename.zip
	`

	var (
		flags   = flag.NewFlagSet(argv[0], flag.ExitOnError)
		fSend   = flags.Bool("send", false, "put data")
		fGet    = flags.Bool("get", false, "get data")
		fConfig = flags.String("config", "", "It performs by the specified directory.")
	)
	flags.StringVar(fConfig, "c", "", "")
	flags.BoolVar(fSend, "s", false, "")
	flags.BoolVar(fGet, "g", false, "")

	flags.Parse(argv[1:])
	if flags.NArg() != 1 || (!*fSend && !*fGet) {
		fmt.Println(usage)
		flags.Usage()
		return errors.Errorf("parse sender command")
	}
	zipFile := flags.Args()[0]
	log.Debugf("set zip %s", zipFile)
	hostName, err := GetHostname()
	if err != nil {
		log.Fatal("get hostname ", err)
	}
	if *fConfig == "" {
		home, err := GetParentAbsPath(os.Args[0], 2)
		if err != nil {
			log.Fatal("get getperf path ", err)
		}
		*fConfig = filepath.Join(home, "getperf.ini")
	}
	configEnv := NewConfigEnvBase(hostName, cmdName, *fConfig)
	home, err := GetParentAbsPath(*fConfig, 1)
	if err != nil {
		log.Fatal("get home ", err)
	}
	config := NewConfig(home, configEnv)
	config.ParseConfigFile(config.ParameterFile)

	// Webサービス設定ファイル読込み
	config.ParseConfigFile(config.SvParamFile)
	if err := SetLogLevel(config.Schedule.LogLevel); err != nil {
		log.Fatal("set log level ", err)
	}
	return config.RunFileTransfer(*fSend, zipFile)
}
