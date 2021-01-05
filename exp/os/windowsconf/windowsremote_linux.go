package windowsconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"

	"github.com/getperf/getperf2/exp/os/windowsconf/winrm-auth-ntlm"
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
	// // Note, username/password pair in the NewClientWithParameters call is ignored
	client, err := winrm.NewClientWithParameters(endpoint, "", "", winrm.DefaultParameters)
	//client, err := winrm.NewClient(endpoint, sv.User, sv.Password)
	if err != nil {
		return HandleError(e.errFile, err, fmt.Sprintf("run %s", sv.Server))
	}

	for _, metric := range append(metrics, e.Metrics...) {
		if metric.Level == -1 || metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" {
			continue
		}
		startTime := time.Now()
		outFile, err := env.OpenServerLog(sv.Server, metric.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		var cmd string
		if metric.Type == "Cmdlet" {
			cmd = winrm.Powershell(metric.Text)
		} else if metric.Type == "Cmd" {
			cmd = fmt.Sprintf("cmd.exe /c \"%s\"", metric.Text)
		} else {
			cmd = metric.Text
		}
		fmt.Fprintf(e.errFile, "run : %s:%s\n", sv.Server, metric.Id)
		if _, err = client.Run(cmd, outFile, e.errFile); err != nil {
			// if strings.Contains(fmt.Sprintf("%s", err), "http error 401") {
			msg := fmt.Sprintf("%s", err)
			if strings.Contains(msg, "error") ||
				strings.Contains(msg, "connection attempt failed") {
				return HandleError(e.errFile, err, fmt.Sprintf("run %s", sv.Server))
			}
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", sv.Server, metric.Id))
		}
		log.Debugf("run %s:%s,elapse %s", sv.Server, metric.Id, time.Since(startTime))
	}

	return nil
}
