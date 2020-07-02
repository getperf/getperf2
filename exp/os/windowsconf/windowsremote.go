package windowsconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dpotapov/winrm-auth-ntlm"
	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/masterzen/winrm"
	log "github.com/sirupsen/logrus"
)

// var (
// 	defaultTimeoutDuration = 300 * time.Second
// 	timeoutKillAfter       = 1 * time.Second
// )

func (e *Windows) RunRemoteServer(ctx context.Context, env *cfg.RunEnv, sv *Server) error {
	log.Infof("server : %v", sv)
	e.datastore = filepath.Join(env.Datastore, sv.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}

	endpoint := winrm.NewEndpoint(sv.Url, 5985, false, false, nil, nil, nil, 0)
	winrm.DefaultParameters.TransportDecorator = func() winrm.Transporter {
		return &winrmntlm.Transport{
			Username: sv.User,
			Password: sv.Password,
		}
	}
	// Note, username/password pair in the NewClientWithParameters call is ignored
	client, err := winrm.NewClientWithParameters(endpoint, "", "", winrm.DefaultParameters)
	if err != nil {
		return HandleError(e.errFile, err, fmt.Sprintf("run %s", sv.Server))
	}

	for _, command := range append(commands, e.Commands...) {
		if command.Level > env.Level {
			continue
		}
		if command.Id == "" {
			continue
		}
		startTime := time.Now()
		outFile, err := env.OpenServerLog(sv.Server, command.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		if _, err = client.Run(command.Text, outFile, e.errFile); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", sv.Server, command.Id))
		}
		log.Infof("run %s:%s,elapse %s", sv.Server, command.Id, time.Since(startTime))
	}

	// 	startTime := time.Now()
	// 	outFile, err := env.OpenLog(command.Id)
	// 	if err != nil {
	// 		return errors.Wrap(err, "prepare windows inventory log")
	// 	}
	// 	defer outFile.Close()
	// 	cmdContext := append(cmdPrefix, command.Text)
	// 	cmd := exec.Command(cmdContext[0], cmdContext[1:]...)
	// 	cmd.Stdout = outFile
	// 	cmd.Stderr = outFile
	// 	tio := &timeout.Timeout{
	// 		Cmd:       cmd,
	// 		Duration:  defaultTimeoutDuration,
	// 		KillAfter: timeoutKillAfter,
	// 	}
	// 	if env.Timeout != 0 {
	// 		tio.Duration = time.Duration(env.Timeout) * time.Second
	// 	}
	// 	exit, err := tio.RunContext(ctx)
	// 	if err != nil {
	// 		return errors.Wrap(err, "run windows inventory process")
	// 	}
	// 	// exit := <-ch
	// 	msg := fmt.Sprintf("%s,RC:%d,Signal:%t,Elapse:%s",
	// 		command.Text,
	// 		exit.GetChildExitCode(), exit.Signaled, time.Since(startTime))
	// 	if exit.GetChildExitCode() != 0 || exit.Signaled {
	// 		log.Error(msg)
	// 	}
	// 	log.Infof("Complete command %s", msg)
	// }
	return nil
}
