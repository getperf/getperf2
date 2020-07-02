package agent

import (
	// "bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var stubcmd = "../testdata/stubcmd"

func init() {
	if runtime.GOOS == "windows" {
		stubcmd = `..\testdata\stubcmd.exe`
	}
	err := exec.Command("go", "build", "-o", stubcmd, "../testdata/stubcmd.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestCheckServiceExist(t *testing.T) {
	config := CreateTempWorkDirConfig()
	defer os.RemoveAll(config.WorkDir)
	config.PidFile = "_pid_getperf"
	config.WriteWorkFileNumber(config.PidFile, 123)
	if ok := config.CheckServiceExist(123); ok == false {
		t.Error("check service exit 1")
	}
	config.WriteWorkFileNumber(config.PidFile, 12345)
	if ok := config.CheckServiceExist(123); ok == true {
		t.Error("check service exit 2")
	}
	os.Remove(config.GetWorkfilePath(config.PidFile))
	if ok := config.CheckServiceExist(1); ok == true {
		t.Error("check service exit 3")
	}
}

func TestCheckProcess(t *testing.T) {
	pid := os.Getpid()
	if CheckProcess(pid, "") == false {
		t.Error("check own process")
	}
	if CheckProcess(pid, "hogehoge") == true {
		t.Error("check process wrong keyword")
	}
	if CheckProcess(-1, "") == true {
		t.Error("check process wrong pid")
	}
}

func TestExecCommandRedirect(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	CmdInfo := &CommandInfo{
		CmdLine: "../testdata/stubcmd -sleep 15s",
		OutPath: filepath.Join(logDir, "out.txt"),
		Timeout: 2,
	}
	err := CmdInfo.ExecCommandRedirect()
	if CmdInfo.Executed == false {
		t.Error("command not execute")
	}
	t.Log("Error : ", err)
	t.Log("Executed : ", CmdInfo.Executed)
	t.Log("Pid : ", CmdInfo.Pid)
	t.Log("ExitCode : ", CmdInfo.ExitCode)
	t.Log("Status : ", CmdInfo.Status)
	t.Log("OutBuf : ", CmdInfo.OutBuf)

	outBuf, _ := ioutil.ReadFile(CmdInfo.OutPath)
	// z := bytes.NewBuffer(outBuf)
	if decodeBytes(outBuf) != "OUT TEST\nERR TEST\n" {
		t.Error("command output")
	}
}

func TestExecCommandIncludeDirectoryPath(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	CmdInfo := &CommandInfo{
		CmdLine: "../testdata/stubcmd",
		OutPath: filepath.Join(logDir, "stub/out.txt"),
		Timeout: 2,
	}
	err := CmdInfo.ExecCommandRedirect()
	if CmdInfo.Executed == false {
		t.Error("command not execute redirect outfile including directory")
	}
	t.Log("ERR : ", err)
	t.Log("RES : ", *CmdInfo)

	outBuf, _ := ioutil.ReadFile(CmdInfo.OutPath)
	// z := bytes.NewBuffer(outBuf)
	if decodeBytes(outBuf) != "OUT TEST\nERR TEST\n" {
		t.Error("command output")
	}
}

func TestExecCommandNoRedirect(t *testing.T) {
	// SetLogLevel(7)
	CmdInfo := &CommandInfo{
		CmdLine: "../testdata/stubcmd -sleep 5s",
		Timeout: 2,
	}
	err := CmdInfo.ExecCommandNoRedirect()
	if CmdInfo.Executed == false {
		t.Error("command not execute")
	}
	t.Log("ERR : ", err)
	t.Log("RES : ", *CmdInfo)
	t.Log("OUT : ", CmdInfo.OutBuf)
	if CmdInfo.OutBuf != "OUT TEST\nERR TEST\n" {
		t.Error("command output")
	}
}
