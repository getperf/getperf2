package linuxconf

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Songmu/timeout"
	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	log "github.com/sirupsen/logrus"
)

var cmdBase = []string{"sh", "-c"}

func (e *Linux) RunLocalCommand(ctx context.Context, metric *Metric) error {
	cmdArgs := append(cmdBase, filepath.FromSlash(metric.Text))

	log.Debug("exec command direct ", cmdArgs)
	args := append([]string{}, cmdArgs...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = metric.stdOut
	cmd.Stderr = metric.stdErr
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	if e.Env.Timeout != 0 {
		tio.Duration = time.Duration(e.Env.Timeout) * time.Second
	}
	exitStatus, err := tio.RunContext(ctx)
	if err != nil {
		fmt.Fprintf(metric.stdErr, "%s:%s\n", metric.Id, err)
	}
	pid := cmd.ProcessState.Pid()
	if err == nil && (exitStatus.IsTimedOut() || exitStatus.Signaled) {
		fmt.Fprintf(metric.stdErr, "timeout, pid %d\n", pid)
	}
	exitCode := exitStatus.GetChildExitCode()
	if exitCode != 0 {
		fmt.Fprintf(metric.stdErr, "exit %d, pid %d\n", exitCode, pid)
	}
	return nil
}

func (e *Linux) RunLocalServer(ctx context.Context, env *cfg.RunEnv, server string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("linux environment only")
	}
	e.datastore = filepath.Join(env.Datastore, server)
	log.Info("collect local server : ", server, e.datastore)

	if err := RemoveAndCreateDir(e.datastore); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	for _, metric := range metrics {
		if metric.Level > env.Level {
			continue
		}
		if metric.Id == "" {
			continue
		}
		startTime := time.Now()
		outFile, err := env.OpenServerLog(server, metric.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		metric.stdOut = outFile
		metric.stdErr = e.errFile
		if err := e.RunLocalCommand(ctx, metric); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", server, metric.Id))
		}
		log.Infof("run %s:%s,elapse %s", server, metric.Id, time.Since(startTime))
	}
	return nil
}
