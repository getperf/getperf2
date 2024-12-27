package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
)

func (config *Config) StopDaemon() error {
	persistentPid, err := config.ReadWorkFileNumber(config.PidFile)
	if err != nil {
		return errors.Wrap(err, "read pid file for stop service")
	}
	if !CheckProcess(persistentPid, "getperf") {
		return fmt.Errorf("stop service failed for can't find pid file process : %d", persistentPid)
	}
	log.Infof("found service pid : %d", persistentPid)
	if err = config.WriteWorkFile(config.ExitFlag, []string{"stop"}); err != nil {
		return errors.Wrap(err, "write exit flag file for stop service")
	}
	log.Infof("Waiting %d sec for shutting down the getperf process",
		DEFAULT_SOAP_TIMEOUT)
	// timeout := DEFAULT_SOAP_TIMEOUT
	timeout := 5
	for timeout > 0 {
		exist, err := config.CheckWorkFile(config.ExitFlag)
		if err != nil {
			return errors.Wrap(err, "check exit flag file for stop service")
		}
		if !exist {
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
		timeout = timeout - 1
		if timeout == 0 {
			KillProcess(persistentPid, "getperf")
		}
	}
	return nil

}

func (c *Config) RunDaemon() {
	persistentPid, err := c.ReadWorkFileNumber(c.PidFile)
	if err == nil {
		if CheckProcess(persistentPid, "getperf") {
			fmt.Printf("start up failed for another pid file process : %d\n", persistentPid)
			return
		}
	}
	agentService := &daemon.Context{
		Umask: 027,
		Args:  os.Args,
	}

	d, err := agentService.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer agentService.Release()

	log.Println("daemon started")

	go Run(context.Background(), os.Args[0:], os.Stdout, os.Stderr)

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}
