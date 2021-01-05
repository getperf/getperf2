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
	ps "github.com/hnakamur/go-powershell"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	defaultTimeoutDuration = 300 * time.Second
	timeoutKillAfter       = 1 * time.Second
)

var powershellCmd = `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`

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

func (e *Windows) writeScript(templateCode string, doc io.Writer, env *cfg.RunEnv) error {
	// tmpl, err := template.ParseFiles("powershell.tpl")
	tmpl, err := template.New("windowsconf").Parse(templateCode)
	if err != nil {
		return errors.Wrap(err, "failed read template")
	}
	var filteredMetrics []*Metric
	for _, metric := range append(e.Metrics) {
		if metric.Level == -1 || metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" {
			continue
		}
		if metric.Type == "Cmd" {
			metric.Text = fmt.Sprintf("cmd.exe /c \"%s\"", metric.Text)
		}
		log.Debugf("add test item %s:%d:%d", metric.Id, metric.Level, env.Level)
		filteredMetrics = append(filteredMetrics, metric)
	}
	if err := tmpl.Execute(doc, filteredMetrics); err != nil {
		return errors.Wrap(err, "failed generate script")
	}
	return nil
}

func (e *Windows) CreateScript(templateCode string, script string, env *cfg.RunEnv) error {
	log.Info("create temporary log dir for test ", env.Datastore)
	// e.ScriptPath = filepath.Join(env.Datastore, "get_windows_inventory.ps1")
	e.ScriptPath = filepath.Join(env.Datastore, script)
	outFile, err := os.OpenFile(e.ScriptPath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed create script")
	}
	defer outFile.Close()
	e.writeScript(templateCode, outFile, env)
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
		if metric.Level == -1 || metric.Level > env.Level {
			continue
		}
		if metric.Id == "" || metric.Text == "" || metric.Type != "Cmd" {
			continue
		}
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
	// if err := e.RunCommands(ctx, env); err != nil {
	// 	return errors.Wrap(err, "run external command")
	// }
	if env.DryRun {
		return e.writeScript(powershellCode, os.Stdout, env)
	}
	startTime := time.Now()
	e.datastore = filepath.Join(env.Datastore, e.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	if err := e.CreateScript(powershellCode, "get_windows_inventory.ps1", env); err != nil {
		return HandleError(e.errFile, err, "prepare script")
	}
	// cmdPowershell := []string{
	// 	// Get-NetConnectionProfileなど、一部コマンドレットの実行で、
	// 	// "プロバイダーによる読み込みエラーです" エラーが発生。
	// 	// 以下の記事を参照し、絶対パスで 64bit 版 PowerShell を指定して
	// 	// も同様のエラーが発生する。原因、調査中。
	// 	// Get-WindowsFeature コマンドレットでも類似の問題発生。

	// 	// https://stackoverflow.com/questions/28156066/how-to-resolve-get-netconnectionprofile-provider-load-failure-on-x86-powershel

	// 	// powershellCmd,
	// 	"powershell",
	// 	e.ScriptPath,
	// 	env.Datastore,
	// }
	shell, err := ps.New()
	if err != nil {
		return HandleError(e.errFile, err, "prepare winrm session")
	}
	defer shell.Exit()

	outFile, err := env.OpenLog("output.log")
	if err != nil {
		return HandleError(e.errFile, err, "prepare shell output log")
	}
	defer outFile.Close()

	cmd := fmt.Sprintf(
		"powershell -ExecutionPolicy RemoteSigned %s %s",
		e.ScriptPath, e.datastore)
	stdout, err := shell.Exec(cmd)
	if err != nil {
		HandleError(e.errFile, err, "exec powershell script")
	}
	outFile.Write([]byte(stdout))

	// cmd := exec.Command(cmdPowershell[0], cmdPowershell[1:]...)
	// cmd.Stdout = outFile
	// cmd.Stderr = errFile
	// tio := &timeout.Timeout{
	// 	Cmd:       cmd,
	// 	Duration:  defaultTimeoutDuration,
	// 	KillAfter: timeoutKillAfter,
	// }
	// if env.Timeout != 0 {
	// 	tio.Duration = time.Duration(env.Timeout) * time.Second
	// }
	// exit, err := tio.RunContext(ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "run windows inventory process")
	// }
	// // exit := <-ch
	// msg := fmt.Sprintf("RC:%d,Signal:%t,Elapse %s",
	// 	exit.GetChildExitCode(), exit.Signaled, time.Since(startTime))
	// if exit.GetChildExitCode() != 0 || exit.Signaled {
	// 	log.Error(msg)
	// 	return errors.New(msg)
	// }
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
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
		if err = e.RunLocalServer(ctx, env); err != nil {
			msg := fmt.Sprintf("run local server '%s'", e.Server)
			HandleError(e.errFile, err, msg)
		}
	}
	for _, sv := range e.Servers {
		log.Info("collect remote server : ", sv.Server)
		if err = e.RunRemoteServer(ctx, env, sv); err != nil {
			msg := fmt.Sprintf("run remote server '%s'", sv.Server)
			HandleError(e.errFile, err, msg)
		}
	}
	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete Windows inventory collection %s", msg)

	return err
}

// func (e *Windows) Run(env *cfg.RunEnv) error {
// 	if runtime.GOOS != "windows" {
// 		return fmt.Errorf("windows powershell environment only")
// 	}
// 	// startTime := time.Now()
// 	if env.DryRun {
// 		return e.writeScript(powershellCode, os.Stdout, env)
// 	}
// 	if err := e.CreateScript(powershellCode, "get_windows_inventory.ps1", env); err != nil {
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
