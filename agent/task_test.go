package agent

import (
	// "bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewTaskJob(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.Linux = 'netstat -s', netstat.txt")
	job := schedule.Collectors["Linux"].Jobs[0]
	taskJob := NewTaskJob(1, job, "", "")
	if taskJob.CmdLine != "netstat -s" {
		t.Error("new task job")
	}
}

func TestNewTaskJobMacro(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/stubcmd > _odir_/out.txt'")
	job := schedule.Collectors["HW"].Jobs[0]
	taskJob := NewTaskJob(1, job, "/hoge/log", "/hoge/script")
	if taskJob.CmdLine != "/hoge/script/stubcmd > /hoge/log/out.txt" {
		t.Error("new task job macro")
	}
}

func TestNewTask(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.Linux = 'netstat -s', netstat.txt")
	task := NewTask(schedule.Collectors["Linux"], "", "")
	taskJob := task.TaskJobs[0]
	if taskJob.Job.Cmd != "netstat -s" {
		t.Error("new task job")
	}
}

func TestTaskJobRun(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")
	task := NewTask(schedule.Collectors["HW"], logDir, "")
	taskJob := task.TaskJobs[0]
	cmdInfo, err := taskJob.Run()
	if cmdInfo.Executed == false || cmdInfo.Pid == 0 || err != nil {
		t.Error("task job run")
	}
}

func TestTaskJobNoRedirectRun(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd > _odir_/out.txt'")
	task := NewTask(schedule.Collectors["HW"], logDir, ".")
	taskJob := task.TaskJobs[0]
	cmdInfo, _ := taskJob.Run()
	if cmdInfo.Executed == false || cmdInfo.Pid == 0 || cmdInfo.OutBuf != "ERR TEST\n" {
		t.Error("task job run")
	}
}

func TestTaskJobPeriodicalRun(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd', out.txt, 2, 5")
	task := NewTask(schedule.Collectors["HW"], logDir, ".")
	taskJob := task.TaskJobs[0]
	taskJob.Timeout = 5
	cmdInfo, _ := taskJob.Run()
	t.Log("CmdInfo : ", *cmdInfo)
	outBuf, _ := ioutil.ReadFile(taskJob.CommandInfo.OutPath)
	// if decodeBytes(bytes.NewBuffer(outBuf)) != "OUT TEST\nERR TEST\nOUT TEST\nERR TEST\nOUT TEST\nERR TEST\n" {
	if decodeBytes(outBuf) != "OUT TEST\nERR TEST\nOUT TEST\nERR TEST\nOUT TEST\nERR TEST\n" {
		t.Error("periodical run")
	}
}

func TestMakeReport(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd > _odir_/out.txt'")
	task := NewTask(schedule.Collectors["HW"], logDir, ".")
	err := task.Run()
	yaml, err := task.MakeReport()
	t.Log("YAML : ", yaml)
	t.Log("Error : ", err)
}

func TestSerialTask(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_MODE.HW = serial")
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd > _odir_/out.txt'")
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd', out2.txt, 2, 5")
	task := NewTask(schedule.Collectors["HW"], logDir, ".")
	err := task.Run()
	yaml, err := task.MakeReport()
	t.Log("YAML : ", yaml)
	t.Log("Error : ", err)
}

func TestConcurrentTask(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_MODE.HW = concurrent")
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd > _odir_/out.txt'")
	schedule.ParseConfigLine("STAT_CMD.HW = '_script_/../testdata/stubcmd', out2.txt, 2, 5")
	task := NewTask(schedule.Collectors["HW"], logDir, ".")
	err := task.Run()
	yaml, err := task.MakeReport()
	t.Log("YAML : ", yaml)
	t.Log("Error : ", err)
}

func TestTaskChDirAndRun(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)
	scriptDir, _ := ioutil.TempDir("", "script")
	defer os.RemoveAll(scriptDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_MODE.HW = serial")
	schedule.ParseConfigLine("STAT_CMD.HW = 'pwd > _odir_/out.txt'")
	task := NewTask(schedule.Collectors["HW"], logDir, scriptDir)
	err := task.Run()
	if err != nil {
		t.Error("task job run", err)
	}
}

func TestTaskReport(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	defer os.RemoveAll(logDir)

	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_MODE.HW = serial")
	schedule.ParseConfigLine("STAT_CMD.HW = \"0123456789 01234567890123456789012345678901234567890123456789012345678901234567890123456 78901234567890123456789012345678901234567890123456789012345678901234567890123 45678901234567890123456789012345678901234567890123456789012345678901234567890 1234567890123456789012345678901234567890123456789012 3456789\"")
	task := NewTask(schedule.Collectors["HW"], logDir, "")
	yaml, err := task.MakeReport()
	t.Log("YAML : ", yaml)
	t.Log("Error : ", err)
}
