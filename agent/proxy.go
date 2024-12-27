package agent

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func (config *Config) SetupHttpProxy() error {
	log.Debugf("setup proxy %v %v", config.Host, config.Schedule.ProxyEnable)
	schedule := config.Schedule
	if schedule.ProxyEnable && schedule.ProxyHost != "" {
		httpProxy := fmt.Sprintf("http://%s:%d", schedule.ProxyHost, schedule.ProxyPort)
		os.Setenv("HTTP_PROXY", httpProxy)
		log.Infof("set proxy %v", httpProxy)
	} else {
		os.Setenv("HTTP_PROXY", "")
		os.Setenv("HTTPS_PROXY", "")
		os.Setenv("http_proxy", "")
		os.Setenv("https_proxy", "")
	}

	return nil
}
