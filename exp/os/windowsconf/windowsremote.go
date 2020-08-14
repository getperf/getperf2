package windowsconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func convCommandLine(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
		"\t", nlcode,
	).Replace(str)
}

func (e *Windows) RunRemoteServer(ctx context.Context, env *cfg.RunEnv, sv *Server) error {
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
		var cmd string
		if command.Type == "Cmdlet" {
			cmd = winrm.Powershell(command.Text)
		} else if command.Type == "Cmd" {
			cmd = fmt.Sprintf("cmd.exe /c \"%s\"", command.Text)
		} else {
			cmd = command.Text
		}
		fmt.Fprintf(e.errFile, "run : %s:%s\n", sv.Server, command.Id)
		if _, err = client.Run(cmd, outFile, e.errFile); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", sv.Server, command.Id))
		}
		log.Debugf("run %s:%s,elapse %s", sv.Server, command.Id, time.Since(startTime))
	}

	return nil
}
