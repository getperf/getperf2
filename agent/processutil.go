package agent

import (
	"strings"

	"github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

func (config *Config) CheckServiceExist(pid int) bool {
	persistentPid, err := config.ReadWorkFileNumber(config.PidFile)
	if err != nil {
		log.Error("check service exist ", err)
		return false
	}
	return (persistentPid == pid)
}

// github.com/mitchellh/go-ps 実装
func CheckProcess(pid int, keyword string) bool {
	pidInfo, err := ps.FindProcess(pid)
	if err != nil {
		log.Error("check process ", err)
		return false
	}
	if pidInfo == nil {
		return false
	}
	if pidInfo.Pid() != pid {
		return false
	}
	if keyword != "" && strings.Index(pidInfo.Executable(), keyword) == -1 {
		return false
	}
	log.Info("found process ", pidInfo)
	return true
}
