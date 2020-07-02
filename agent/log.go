package agent

import (
	"fmt"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func SetLog(home string) error {
	logPath := filepath.Join(home, "_log", "getperf_log.%Y%m%d")
	rl, err := rotatelogs.New(logPath)
	if err != nil {
		return fmt.Errorf("set log %s", err)
	}
	log.SetOutput(rl)
	return nil
}

// ; Log level. None 0, FATAL 1, CRIT 2, ERR 3, WARN 4, NOTICE 5, INFO 6, DBG 7
// LOG_LEVEL = 5
func SetLogLevel(level int) error {
	switch level {
	case 0, 1:
		log.SetLevel(log.FatalLevel)
	case 2, 3:
		log.SetLevel(log.ErrorLevel)
	case 4, 5:
		log.SetLevel(log.WarnLevel)
	case 6:
		log.SetLevel(log.InfoLevel)
	case 7:
		log.SetLevel(log.DebugLevel)
	default:
		return fmt.Errorf("unkown log level %d", level)
	}
	return nil
}
