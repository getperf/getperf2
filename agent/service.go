package agent

import (
	"context"
	"os"
	"time"

	"github.com/kardianos/service"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type windows_program struct{}

func (p *windows_program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *windows_program) run() {
	go Run(context.Background(), os.Args[0:], os.Stdout, os.Stderr)
}

func (p *windows_program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 8)
	return nil
}

func (config *Config) WindowsServiceControl(command string) error {
	log.Infof("run windows service control %s", command)
	svcConfig := &service.Config{
		Name:        "Getperf2",
		DisplayName: "Getperf Cacti Agent2",
		Description: "Getperf Cacti Agent coded in Go language.",
	}

	prg := &windows_program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := service.Control(s, command); err != nil {
		return errors.Wrap(err, "run service control command")
	}
	return nil
}
