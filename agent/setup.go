package agent

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func checkServiceUrlFormt(serviceUrl string) error {
	if _, err := url.ParseRequestURI(serviceUrl); err != nil {
		return errors.Wrap(err, "check url")
	}
	return nil
}

func (config *Config) RunSetupConsole(siteKey, password, url string) error {
	schedule := config.Schedule
	if siteKey == "" {
		siteKey = schedule.SiteKey
	}
	var err error
	schedule.SiteKey, err = common.ReadLine(Translate("ja", "Enter site key"), siteKey)
	if err != nil {
		return errors.Wrap(err, "read line siteKey")
	}
	schedule.Password, err = common.ReadLine(Translate("ja", "Enter password"), password)
	if err != nil {
		return errors.Wrap(err, "read line password")
	}
	if url != "" {
		if err := checkServiceUrlFormt(url); err != nil {
			return errors.Wrap(err, "check --url option")
		}
		config.Schedule.UrlCM = url + "/axis2/services/GetperfService"
	}
	if config.Schedule.UrlCM == "" {
		return errors.Wrap(err, "URL_CM not found")
	}
	log.Debugf("Setup: site_key=%s, pass=%s, url=%s",
		schedule.SiteKey, schedule.Password, schedule.UrlCM)
	return nil
}

func (config *Config) RunSetup(siteKey, password, urlCM string) error {
	// log.Info("RUN SETUP")
	log.Info("run setup ", VersionMessage())
	persistentPid, err := config.ReadWorkFileNumber(config.PidFile)
	if err != nil {
		log.Infof("read pid file for stop service : %s", err)
		// return errors.Wrapf(err, "read pid file for stop service %d", persistentPid)
	}
	if CheckProcess(persistentPid, "getperf") {
		return fmt.Errorf(
			"Process pid=%d is running. Enter 'getperfctl stop' command, if you stop Agent.",
			persistentPid,
		)
	}
	config.InitAgent()
	config.ParseConfigFile(config.ParameterFile)

	if err = config.LoadLicense(); err != nil {
		log.Info("Initialize SSL license file")
	}
	if err = config.AuthLicense(0); err != nil {
		log.Info("Login failed. Please enter the correct ID or Password or republish your ID on the portal site.")
	}

	if err = config.RunSetupConsole(siteKey, password, urlCM); err != nil {
		return errors.Wrap(err, "invalid setup parameter")
	}

	if err := config.SetupHttpProxy(); err != nil {
		log.Fatal("set http proxy ", err)
	}

	schedule := config.Schedule
	u, err := url.ParseRequestURI(schedule.UrlCM)
	if err != nil {
		return errors.Wrap(err, "parse URL_CM")
	}
	portNum, _ := strconv.Atoi(u.Port())
	log.Infof("set web service host: %v, port: %v", u.Hostname(), portNum)
	soapSender, err := NewSoapSender(u.Hostname(), portNum)
	if err != nil {
		return errors.Wrap(err, "initialize soap sender")
	}
	_, err = soapSender.MakeTransportWithServerAuthSSL(config)
	if err != nil {
		return errors.Wrap(err, "set client cert config")
	}

	sslPath := filepath.Join(config.WorkCommonDir, "sslconf.zip")
	log.Infof("download network config %s", sslPath)
	result, err := soapSender.RegistAgent(schedule.SiteKey, config.Host, schedule.Password, sslPath)
	if err != nil || result != "OK" {
		return errors.Wrap(err, "regist agent")
	}
	log.Infof("unzip network config %v", sslPath)
	return Unzip(sslPath, config.Home)

}
