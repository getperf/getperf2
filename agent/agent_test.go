package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var soapcmdpath = "../testdata/getperfsoap"
var soapcmdfile = "getperfsoap"
var testHost = "centos80"

func init() {
	if runtime.GOOS == "windows" {
		soapcmdpath = `..\testdata\getperfsoap.exe`
		soapcmdfile = `getperfsoap.exe`
	}
	err := exec.Command("go", "build", "-o", soapcmdpath, "../testdata/getperfsoap.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestCheckExitFile(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	c := NewConfig(home, NewConfigEnv())
	c.InitAgent()
	if status, _ := c.CheckExitFile(); status != "" {
		t.Error("check no exit file")
	}
	c.WriteWorkFile(c.ExitFlag, []string{"stop"})
	if status, _ := c.CheckExitFile(); status != "stop" {
		t.Error("check exit file")
	}
	c.WriteWorkFile(c.ExitFlag, []string{})
	if status, _ := c.CheckExitFile(); status != "" {
		t.Error("check exit file2")
	}
}

func TestCheckHostname(t *testing.T) {
	c := NewConfig("ptune_home", NewConfigEnv())
	if c.CheckHostname("moi") == false {
		t.Error("check normal hostname")
	}
	if c.CheckHostname("wrong host name") == true {
		t.Error("check wrong hostname")
	}
	if c.CheckHostname("wrong/host/name") == true {
		t.Error("check wrong separator hostname")
	}
	if c.CheckHostname("wrong\\host\\name") == true {
		t.Error("check wrong separator hostname2")
	}
}

func TestCheckHAStatus(t *testing.T) {
	c := NewConfig("ptune_home", NewConfigEnv())
	schedule := NewSchedule()
	haCmd := stubcmd + " -echo test1"
	schedule.ParseConfigLine("HANODE_ENABLE = true")
	schedule.ParseConfigLine("HANODE_CMD = " + haCmd)
	c.Schedule = schedule
	if err := c.CheckHAStatus(); err != nil && c.ServiceName == "test1" {
		t.Error("check ha status normal")
	}
}

func TestAuthLicense(t *testing.T) {
	c := NewConfig("ptune_home", NewConfigEnv())
	schedule := NewSchedule()
	schedule.ParseConfigFile("../testdata/ptune/network/License.txt")
	c.Host = testHost
	c.Schedule = schedule
	if err := c.AuthLicense(0); err != nil {
		t.Error("auth license normal", err)
	}
}

func TestUnzipSSLConfig(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	c := NewConfig(home, NewConfigEnv())
	c.InitAgent()
	sslPath := filepath.Join(c.WorkCommonDir, "sslconf.zip")
	CopyFile("../testdata/ptune/_wk/sslconf.zip", sslPath)
	t.Log(sslPath)
	if err := c.UnzipSSLConf(); err != nil {
		t.Error("unzip ssl config ", err)
	}
}

// * アーカイバ移行 []
func TestArchiveData(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	// defer os.RemoveAll(home)
	t.Log(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	config.ParseConfigLine("STAT_MODE.HW = serial")
	config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")
	datastore, err := config.NewDatastoreCurrent("HW")
	collector := config.GetCollector("HW")
	task := NewTask(collector, datastore.AbsDir(), "")
	err = task.Run()
	if err != nil {
		t.Error("run task")
	}
	err = config.ArchiveData(task, datastore)
	if err != nil {
		t.Error("archive")
	}
}

func TestTruncateBackupData(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	// defer os.RemoveAll(home)
	t.Log(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	testFiles := []string{
		"arc_centos80__HW_20200528_132000.zip",
		"arc_centos80__HW_20200528_132500.zip",
		"arc_centos80__HW_20200528_133000.zip",
		"arc_centos80__HW_20200528_133500.zip",
		"arc_centos80__VMWare_20200528_133000.zip",
		"arc_centos80__VMWare_20200528_133500.zip",
	}
	for _, testFile := range testFiles {
		src := filepath.Join("../testdata/ptune/backup/", testFiles[0])
		dest := filepath.Join(config.BackupDir, testFile)
		CopyFile(src, dest)
	}

	config.ParseConfigLine("STAT_MODE.HW = serial")
	config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")
	datastore, err := config.NewDatastoreCurrent("HW")
	collector := config.GetCollector("HW")
	task := NewTask(collector, datastore.AbsDir(), "")
	task.Run()
	config.ArchiveData(task, datastore)
	err = config.TruncateBackupData(task, datastore)
	if err != nil {
		t.Error("archive")
	}
}

// * Typeperf テスト
func TestTypeperf(t *testing.T) {
	// home, _ := ioutil.TempDir("", "ptune")
	// defer os.RemoveAll(home)
	config := NewConfig("../testdata/ptune_win/", NewConfigEnv())
	config.InitAgent()
	config.ParseConfigLine("STAT_MODE.HW = serial")
	config.ParseConfigLine("STAT_INTERVAL.HW = 3600")
	// 上書きしますか？の確認メッセージ発生、CPU負荷上昇。タイムアウトは有効に動作
	// config.ParseConfigLine("STAT_CMD.HW = 'typeperf.exe -cf testdata\\ptune_win\\script\\PerfMon\\ProcessorMemory.txt -si  5 -sc 61 -f CSV -o \"_odir_\\ProcessorMemory.csv\"'")
	config.ParseConfigLine("STAT_CMD.HW = 'typeperf.exe -cf ..\\testdata\\ptune_win\\script\\PerfMon\\ProcessorMemory.txt -si  5 -sc 61 -f CSV -o \"_odir_\\ProcessorMemory.csv\"'")
	datastore, err := config.NewDatastoreCurrent("HW")
	collector := config.GetCollector("HW")
	task := NewTask(collector, datastore.AbsDir(), "")
	err = task.Run()
	if err != nil {
		t.Error("run task")
	}
}

// * ライセンス管理移行 []

func TestCheckLicenseOK(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	config.ParseConfigFile("../testdata/ptune/network/License.txt")
	config.Host = testHost
	if err := config.CheckLicense(0); err != nil {
		t.Error("check license normal ", err)
	}
}

func TestCheckDiskUtil(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	config.ParseConfigLine("DISK_CAPACITY = 0")
	disk, _ := CheckDiskFree(".")
	t.Log(config.Schedule.DiskCapacity)
	diskUtil := 100.0 * disk.Free / disk.All
	if config.Schedule.DiskCapacity > int(diskUtil) {
		t.Error("disk util 1")
	}
	if err := config.CheckDiskUtil(); err != nil {
		t.Error(err)
	}
}

// * SOAPコマンド管理 []
func TestExecSOAPCommandPM(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	// defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	t.Log("Home:", home)
	t.Log("SOAP:", soapcmdpath)
	os.MkdirAll(config.BinDir, 0777)
	soapcmddest := filepath.Join(config.BinDir, soapcmdfile)
	CopyFile(soapcmdpath, soapcmddest)
	os.Chmod(soapcmddest, 0777)
	cmdLine := fmt.Sprintf("%s -h", soapcmddest)
	t.Log("Exec ", cmdLine)
	cmdInfo := &CommandInfo{
		CmdLine: cmdLine,
		Timeout: config.Schedule.SoapTimeout,
	}
	err := cmdInfo.ExecCommandNoRedirect()
	if err != nil {
		t.Error("exec getperfsoap")
	}
	if cmdInfo.Executed == false {
		t.Error("exec getperfsoap failed")
	}
	if len(cmdInfo.OutBuf) == 0 {
		t.Error("exec getperfsoap failed")
	}

}
