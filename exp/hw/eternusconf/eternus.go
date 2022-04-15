package eternusconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/common/sshx"
	"github.com/getperf/getperf2/common/telnetx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (e *Eternus) RunRemoteServer(ctx context.Context, env *cfg.RunEnv) error {
	log.Info("collect remote server : ", e.Server)
	e.datastore = filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	// client, err := sshConnect(e.Url, e.User, e.Password, e.SshKeyPath)
	client, err := sshx.SshConnect(e.Url, e.User, e.Password, e.SshKeyPath)
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
		outFile, err := env.OpenServerLog(e.Server, metric.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		if err := sshx.RunCommand(outFile, e.errFile, client, metric.Type, metric.Text); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", e.Server, metric.Id))
		}
		log.Infof("run %s:%s,elapse %s", e.Server, metric.Id, time.Since(startTime))
	}
	return nil
}

func (e *Eternus) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare Eternus inventory error")
	}
	defer errFile.Close()
	e.errFile = errFile
	e.Env = env
	fmt.Println("URL:", e.Url)
	telnetUrl, err := telnetx.ParseUrl(e.Url)
	fmt.Println("URL2:", telnetUrl)
	if err != nil {
		return HandleErrorWithAlert(e.errFile, err, "check url")
	}
	if telnetUrl != "" {
		err = e.RunRemoteServerTelnet(ctx, env)
	} else {
		err = e.RunRemoteServer(ctx, env)
	}
	if err != nil {
		msg := fmt.Sprintf("run remote server '%s'", e.Server)
		HandleErrorWithAlert(e.errFile, err, msg)
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete Eternus inventory collection %s", msg)

	return err
}
