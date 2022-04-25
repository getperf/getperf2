package aixconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/common/sshx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)


func (e *AIX) RunRemoteServer(ctx context.Context, env *cfg.RunEnv, sv *Server) error {
	log.Info("collect remote server : ", sv.Server)
	e.remoteServer = sv.Server
	e.datastore = filepath.Join(env.Datastore, sv.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	// client, err := sshConnect(sv.Url, sv.User, sv.Password, sv.SshKeyPath)

	client, err := sshx.SshConnect(sv.Url, sv.User, sv.Password, sv.SshKeyPath)

	if err != nil {
		return HandleError(e.errFile, err, "connect remote server")
	}
	defer client.Close()
	for _, metric := range e.Metrics {
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
		if err := sshx.RunCommand(outFile, e.errFile, client, metric.Type, metric.Text); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", sv.Server, metric.Id))
		}
		log.Infof("run %s:%s,elapse %s", sv.Server, metric.Id, time.Since(startTime))
	}
	return nil
}

func (e *AIX) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare AIX inventory error")
	}
	defer errFile.Close()
	e.errFile = errFile
	e.Env = env

	if e.LocalExec == true {
		log.Info("collect local server : ", e.LocalExec)
		if err = e.RunLocalServer(ctx, env, e.Server); err != nil {
			msg := fmt.Sprintf("run local server '%s'", e.Server)
			HandleErrorWithAlert(e.errFile, err, msg)
		}
	}
	for _, sv := range e.Servers {
		if err = e.RunRemoteServer(ctx, env, sv); err != nil {
			msg := fmt.Sprintf("run remote server '%s'", sv.Server)
			HandleErrorWithAlert(e.errFile, err, msg)
		}
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete AIX inventory collection %s", msg)

	return err
}
