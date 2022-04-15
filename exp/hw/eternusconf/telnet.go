package eternusconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/common/telnetx"
	log "github.com/sirupsen/logrus"
)

func (e *Eternus) RunRemoteServerTelnet(ctx context.Context, env *cfg.RunEnv) error {
	log.Info("collect remote server : ", e.Server)
	e.datastore = filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	telnetUrl, _ := telnetx.ParseUrl(e.Url)
	// client, err := sshConnect(e.Url, e.User, e.Password, e.SshKeyPath)
	client, err := telnetx.TelnetConnect(telnetUrl, e.User, e.Password)
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
		var data string
		if metric.Type == "Cmd" || metric.Type == "" {
			data, err = client.ExecCommand(metric.Text)
		} else {
			data, err = client.ExecScript(metric.Text)
		}
		if err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", e.Server, metric.Id))
		}
		outFile.Write([]byte(data))
		log.Infof("run %s:%s,elapse %s", e.Server, metric.Id, time.Since(startTime))
	}
	return nil
}
