package main

import (
	"time"

	"github.com/jefurry/logrus/hooks/rotatelog"
	rlog "github.com/jefurry/logrus/hooks/rotatelog"
	log "github.com/sirupsen/logrus"
)

func main() {
	log := log.New()
	hook, err := rlog.NewHook("./access_log.%Y%m%d",
		//rotatelog.WithLinkName("./access_log"),
		rotatelog.WithMaxAge(24*time.Hour),
		rotatelog.WithRotationTime(time.Hour),
		rotatelog.WithClock(rotatelog.UTC))

	if err != nil {
		log.Hooks.Add(hook)
	}

}
