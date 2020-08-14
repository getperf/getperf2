package windowsconf

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/Songmu/timeout"
	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	defaultTimeoutDuration = 300 * time.Second
	timeoutKillAfter       = 1 * time.Second
)

var powershellCode = `# Windows invetory collecting script

Param(
    [string]$log_dir
)
$log_dir = Convert-Path $log_dir

{{range $i, $v := .}}
echo TestId::{{$v.Id}}
$log_path = Join-Path $log_dir "{{$v.Id}}"
Invoke-Command  -ScriptBlock { {{$v.Text}} } | Out-File $log_path -Encoding UTF8

{{end}}
`

func (e *Windows) writeScript(doc io.Writer, env *cfg.RunEnv) error {
	// tmpl, err := template.ParseFiles("powershell.tpl")
	tmpl, err := template.New("windowsconf").Parse(powershellCode)
	if err != nil {
		return errors.Wrap(err, "failed read template")
	}
	var filteredCommands []*Command
	for _, command := range append(commands) {
		if command.Level > env.Level {
			continue
		}
		if command.Id == "" {
			continue
		}
		log.Debugf("add test item %s:%d:%d", command.Id, command.Level, env.Level)
		filteredCommands = append(filteredCommands, command)
	}
	if err := tmpl.Execute(doc, filteredCommands); err != nil {
		return errors.Wrap(err, "failed generate script")
	}
	return nil
}

func (e *Windows) CreateScript(env *cfg.RunEnv) error {
	log.Info("create temporary log dir for test ", env.Datastore)
	e.ScriptPath = filepath.Join(env.Datastore, "get_windows_inventory.ps1")
	outFile, err := os.OpenFile(e.ScriptPath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed create script")
	}
	defer outFile.Close()
	e.writeScript(outFile, env)
	return nil
}

func (e *Windows) RunCommands(ctx context.Context, env *cfg.RunEnv) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("windows powershell environment only")
	}
	cmdPrefix := []string{
		"cmd",
		"/c",
	}

	for _, metric := range e.Metrics {
		startTime := time.Now()
		outFile, err := env.OpenLog(metric.Id)
		if err != nil {
			return errors.Wrap(err, "prepare windows inventory log")
		}
		defer outFile.Close()
		cmdContext := append(cmdPrefix, metric.Text)
		cmd := exec.Command(cmdContext[0], cmdContext[1:]...)
		cmd.Stdout = outFile
		cmd.Stderr = outFile
		tio := &timeout.Timeout{
			Cmd:       cmd,
			Duration:  defaultTimeoutDuration,
			KillAfter: timeoutKillAfter,
		}
		if env.Timeout != 0 {
			tio.Duration = time.Duration(env.Timeout) * time.Second
		}
		exit, err := tio.RunContext(ctx)
		if err != nil {
			return errors.Wrap(err, "run windows inventory process")
		}
		// exit := <-ch
		msg := fmt.Sprintf("%s,RC:%d,Signal:%t,Elapse:%s",
			metric.Text,
			exit.GetChildExitCode(), exit.Signaled, time.Since(startTime))
		if exit.GetChildExitCode() != 0 || exit.Signaled {
			log.Error(msg)
		}
		log.Infof("Complete command %s", msg)
	}
	return nil
}

func (e *Windows) RunLocalServer(ctx context.Context, env *cfg.RunEnv) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("windows powershell environment only")
	}
	if err := e.RunCommands(ctx, env); err != nil {
		return errors.Wrap(err, "run external command")
	}
	startTime := time.Now()
	if env.DryRun {
		return e.writeScript(os.Stdout, env)
	}
	if err := e.CreateScript(env); err != nil {
		return errors.Wrap(err, "prepare windows inventory script")
	}
	cmdPowershell := []string{
		"powershell",
		e.ScriptPath,
		env.Datastore,
	}
	outFile, err := env.OpenLog("output.log")
	if err != nil {
		return errors.Wrap(err, "prepare windows inventory log")
	}
	defer outFile.Close()

	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare windows inventory error")
	}
	defer errFile.Close()

	cmd := exec.Command(cmdPowershell[0], cmdPowershell[1:]...)
	cmd.Stdout = outFile
	cmd.Stderr = errFile
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	if env.Timeout != 0 {
		tio.Duration = time.Duration(env.Timeout) * time.Second
	}
	exit, err := tio.RunContext(ctx)
	if err != nil {
		return errors.Wrap(err, "run windows inventory process")
	}
	// exit := <-ch
	msg := fmt.Sprintf("RC:%d,Signal:%t,Elapse %s",
		exit.GetChildExitCode(), exit.Signaled, time.Since(startTime))
	if exit.GetChildExitCode() != 0 || exit.Signaled {
		log.Error(msg)
		return errors.New(msg)
	}
	log.Infof("Complete Windows inventory collection %s", msg)

	return nil
}

func (e *Windows) Run(ctx context.Context, env *cfg.RunEnv) error {
	startTime := time.Now()
	errFile, err := env.OpenLog("error.log")
	if err != nil {
		return errors.Wrap(err, "prepare windows inventory error")
	}
	defer errFile.Close()
	e.errFile = errFile

	if e.LocalExec == true {
		log.Info("collect local server : ", e.LocalExec)
		if err := e.RunLocalServer(ctx, env); err != nil {
			msg := fmt.Sprintf("run local server '%s'", e.Server)
			HandleError(e.errFile, err, msg)
		}
	}
	for _, sv := range e.Servers {
		log.Info("collect remote server : ", sv.Server)
		if err := e.RunRemoteServer(ctx, env, sv); err != nil {
			msg := fmt.Sprintf("run remote server '%s'", sv.Server)
			HandleError(e.errFile, err, msg)
		}
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete Windows inventory collection %s", msg)

	return nil
}

// func (e *Windows) Run(env *cfg.RunEnv) error {
// 	if runtime.GOOS != "windows" {
// 		return fmt.Errorf("windows powershell environment only")
// 	}
// 	// startTime := time.Now()
// 	if env.DryRun {
// 		return e.writeScript(os.Stdout, env)
// 	}
// 	if err := e.CreateScript(env); err != nil {
// 		return errors.Wrap(err, "prepare windows inventory script")
// 	}
// 	cmdPowershell := []string{
// 		"powershell",
// 		e.ScriptPath,
// 		env.Datastore,
// 	}
// 	outFile, err := OpenLog(env, "output.log")
// 	if err != nil {
// 		return errors.Wrap(err, "prepare windows inventory log")
// 	}
// 	defer outFile.Close()

// 	errFile, err := OpenLog(env, "error.log")
// 	if err != nil {
// 		return errors.Wrap(err, "prepare windows inventory error")
// 	}
// 	defer errFile.Close()

// 	log.Info(cmdPowershell)
// 	cmd := exec.metric(cmdPowershell[0], cmdPowershell[1:]...)
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return errors.Wrap(err, "create windows inventory Metric pipe")
// 	}
// 	startTime := time.Now()
// 	testId := ""
// 	intervalTime := startTime
// 	cmd.Start()
// 	scanner := bufio.NewScanner(stdout)
// 	for scanner.Scan() {
// 		if testId != "" {
// 			log.Info(testId, ", Elapse: ", time.Since(intervalTime))
// 		}
// 		testId = scanner.Text()
// 		intervalTime = time.Now()
// 	}
// 	if testId != "" {
// 		log.Info(testId, ", Elapse: ", time.Since(intervalTime))
// 	}
// 	cmd.Wait()
// 	// tio := &timeout.Timeout{
// 	// 	Cmd:       cmd,
// 	// 	Duration:  defaultTimeoutDuration,
// 	// 	KillAfter: timeoutKillAfter,
// 	// }
// 	// exitstatus, stdout, stderr, err := tio.Run()

// 	log.Infof("finish windows inventory script, elapse [%s]", time.Since(startTime))

// 	return nil
// }
