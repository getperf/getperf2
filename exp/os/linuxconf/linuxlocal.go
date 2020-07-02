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

func (e *Linux) RunLocalCommand(ctx context.Context, command *Command) error {
	cmdArgs := append(cmdBase, filepath.FromSlash(command.Text))

	log.Debug("exec command direct ", cmdArgs)
	args := append([]string{}, cmdArgs...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = command.stdOut
	cmd.Stderr = command.stdErr
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
		fmt.Fprintf(command.stdErr, "%s:%s\n", command.Id, err)
	}
	pid := cmd.ProcessState.Pid()
	if err == nil && (exitStatus.IsTimedOut() || exitStatus.Signaled) {
		fmt.Fprintf(command.stdErr, "timeout, pid %d\n", pid)
	}
	exitCode := exitStatus.GetChildExitCode()
	if exitCode != 0 {
		fmt.Fprintf(command.stdErr, "exit %d, pid %d\n", exitCode, pid)
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
	for _, command := range commands {
		if command.Level > env.Level {
			continue
		}
		if command.Id == "" {
			continue
		}
		startTime := time.Now()
		outFile, err := env.OpenServerLog(server, command.Id)
		if err != nil {
			return HandleError(e.errFile, err, "prepare inventory log")
		}
		defer outFile.Close()
		command.stdOut = outFile
		command.stdErr = e.errFile
		if err := e.RunLocalCommand(ctx, command); err != nil {
			HandleError(e.errFile, err, fmt.Sprintf("run %s:%s", server, command.Id))
		}
		log.Infof("run %s:%s,elapse %s", server, command.Id, time.Since(startTime))
	}
	return nil
}
