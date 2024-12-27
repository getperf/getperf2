package agent

import (
	"os"
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

func CheckProcessExist(pid int, keyword string) bool {
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
	log.Debug("found process ", pidInfo)
	return true
}

// github.com/mitchellh/go-ps 実装
func CheckProcess(pid int, keyword string) bool {
	return CheckProcessExist(pid, keyword)
	// pidInfo, err := ps.FindProcess(pid)
	// if err != nil {
	// 	log.Error("check process ", err)
	// 	return false
	// }
	// if pidInfo == nil {
	// 	return false
	// }
	// if pidInfo.Pid() != pid {
	// 	return false
	// }
	// if keyword != "" && strings.Index(pidInfo.Executable(), keyword) == -1 {
	// 	return false
	// }
	// log.Debug("found process ", pidInfo)
	// return true
}

func KillProcess(pid int, keyword string) bool {
	if CheckProcessExist(pid, keyword) {
		log.Info("kill process ", pid)
		process, _ := os.FindProcess(pid)
		err := process.Kill()
		if err != nil {
			log.Error("kill process", err)
			return false
		}
	}
	// pidInfo, err := ps.FindProcess(pid)
	// if err != nil {
	// 	log.Error("check process ", err)
	// 	return false
	// }
	// if pidInfo == nil {
	// 	return false
	// }
	// if pidInfo.Pid() != pid {
	// 	return false
	// }
	// if keyword != "" && strings.Index(pidInfo.Executable(), keyword) == -1 {
	// 	return false
	// }
	// log.Info("kill process ", pidInfo)
	// process, err := os.FindProcess(pid)
	// err = process.Kill()
	// if err != nil {
	// 	log.Error("kill process", err)
	// 	return false
	// }
	return true
}
