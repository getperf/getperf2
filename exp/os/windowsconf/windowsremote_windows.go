package windowsconf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	ps "github.com/hnakamur/go-powershell"
	"github.com/labstack/gommon/log"
	// . "github.com/getperf/getperf2/common"
)

// var (
// 	defaultTimeoutDuration = 300 * time.Second
// 	timeoutKillAfter       = 1 * time.Second
// )

var powershellCodeRemote = `# Windows remote invetory collecting script

Param(
    [string]$log_dir
  , [string]$ip
  , [string]$server
  , [string]$user
  , [string]$password
)
$log_dir = Convert-Path $log_dir
$secure   = ConvertTo-SecureString $password -asplaintext -force
$cred     = New-Object System.Management.Automation.PsCredential $user, $secure

$ErrorActionPreference = "Stop"
$session = $null
try {
    $script:session  = New-PSSession $ip -Credential $cred
} catch [Exception] {
    Write-Error "$error"
    exit 1
}
$ErrorActionPreference = "Continue"

$log_dir = Convert-Path $log_dir

{{range $i, $v := .}}
echo TestId::{{$v.Id}}
$log_path = Join-Path $log_dir "{{$v.Id}}"
Invoke-Command -Session $session  -ScriptBlock { {{$v.Text}} } | Out-File $log_path -Encoding UTF8
{{end}}

Remove-PSSession $session
`

func convCommandLine(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
		"\t", nlcode,
	).Replace(str)
}

func (e *Windows) RunRemoteServer(ctx context.Context, env *cfg.RunEnv, sv *Server) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("windows powershell environment only")
	}
	startTime := time.Now()
	if env.DryRun {
		return e.writeScript(powershellCodeRemote, os.Stdout, env)
	}
	e.datastore = filepath.Join(env.Datastore, sv.Server)
	if err := os.MkdirAll(e.datastore, 0755); err != nil {
		return HandleError(e.errFile, err, "create log directory")
	}
	if err := e.CreateScript(powershellCodeRemote,
		"get_windows_inventory_remote.ps1",
		env); err != nil {
		return HandleError(e.errFile, err, "prepare script")
	}
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
		"powershell -ExecutionPolicy RemoteSigned %s %s %s %s %s %s",
		e.ScriptPath, e.datastore, sv.Url, sv.Server, sv.User, sv.Password)
	stdout, err := shell.Exec(cmd)
	// stdout, err := shell.Exec("powershell -ExecutionPolicy RemoteSigned .\\get_windows_spec.ps1 .\\log 192.168.0.20 w2019 administrator P@ssw0rd20A")
	if err != nil {
		HandleError(e.errFile, err, "exec powershell script")
	}
	outFile.Write([]byte(stdout))

	msg := fmt.Sprintf("Elapse %s", time.Since(startTime))
	log.Infof("Complete Windows inventory collection %s", msg)

	return nil
}
