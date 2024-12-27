package agent

import (
	"strings"
	"testing"
)

func TestParamLineTrim(t *testing.T) {
	if line := paramLineTrim("; comment "); line != "" {
		t.Error("trim comment")
	}
	if line := paramLineTrim("\t   ; comment "); line != "" {
		t.Error("trim comment with space")
	}
}

func TestParseConfigLine(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("DISK_CAPACITY=100")
	schedule.ParseConfigLine("DISK_CAPACITY=hoge")
	if schedule.DiskCapacity != 100 {
		t.Error("parse number")
	}
}

func TestParseConfigLineBool(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("DEBUG_CONSOLE = true")
	schedule.ParseConfigLine("DEBUG_CONSOLE=")
	schedule.ParseConfigLine("DEBUG_CONSOLE = false")
	schedule.ParseConfigLine("DEBUG_CONSOLE=hoge")
	if schedule.DebugConsole != false {
		t.Error("parse bool")
	}
}

func TestParseConfigLineCommand(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("HANODE_CMD = hastat.pl -s service1")
	schedule.ParseConfigLine("POST_SOAP_CMD_TYPE = Internal")
	if schedule.HanodeCmd != "hastat.pl -s service1" || schedule.PostSoapCmdType != "Internal" {
		t.Error("parse command")
	}
}

func TestParseCollectorLineCommand(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_ENABLE.Linux = true")
	if c := schedule.FindOrCreateCollector("Linux"); c.StatEnable != true {
		t.Error("parse collector enable")
	}
	schedule.ParseConfigLine("STAT_ENABLE.01_Linux = true")
	if c := schedule.FindOrCreateCollector("01_Linux"); c.StatEnable != true {
		t.Error("parse collector special charactor")
	}
	schedule.ParseConfigLine("BUILD.Linux = 123")
	if c := schedule.FindOrCreateCollector("Linux"); c.Build != 123 {
		t.Error("parse collector enable")
	}
	schedule.ParseConfigLine("STAT_INTERVAL.Linux = 3600")
	if c := schedule.FindOrCreateCollector("Linux"); c.StatInterval != 3600 {
		t.Error("parse collector interval")
	}
	schedule.ParseConfigLine("STAT_TIMEOUT.Linux = 600")
	if c := schedule.FindOrCreateCollector("Linux"); c.StatTimeout != 600 {
		t.Error("parse collector timeout")
	}
	schedule.ParseConfigLine("STAT_MODE.Linux = serial")
	if c := schedule.FindOrCreateCollector("Linux"); c.StatMode != "serial" {
		t.Error("parse collector mode")
	}
}

func TestParseWorkerLineCommand(t *testing.T) {
	if job := stringToJob("\"vmstat 3 3 > _odir_/vmstat.txt\""); job.Cmd != "vmstat 3 3 > _odir_/vmstat.txt" {
		t.Error("parse job 1 ", job.Cmd)
	}
	if job := stringToJob("vmstat 3 3 > _odir_/vmstat.txt"); job != nil {
		t.Error("parse invalid job")
	}
	if job := stringToJob("'vmstat 3 3 > _odir_/vmstat.txt\""); job != nil {
		t.Error("parse invalid job 2")
	}
	if job := stringToJob("'netstat -s', netstat.txt"); job.Ofile != "netstat.txt" {
		t.Error("parse invalid job 3 : ", job.Ofile)
	}
	if job := stringToJob("\"netstat -s\"\t, netstat.txt"); job.Ofile != "netstat.txt" {
		t.Errorf("parse invalid job 4 : %v", job)
	}
	if job := stringToJob("\"netstat -s\"\t, netstat.txt, 30, 10"); job.Step != 10 {
		t.Error("parse invalid job 5 : ", job.Step)
	}
}

func TestParsePsCommand(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.Process = \"/bin/ps -eo pid,ppid,group,user,time,vsz,rss,command\", psutil.txt, 30, 11")
	job := schedule.Collectors["Process"].Jobs[0]
	t.Log("JOB:", job)
}

func TestParserPsRegex(t *testing.T) {
	// "文字列"|'文字列',... を解析、
	tests := []struct {
		command  string
		interval int
		count    int
	}{
		{command: "\"/bin/ps 'test' -eo pid,ppid,group,user,time,vsz,rss,command\""},
		{command: "'/bin/ps \"test\" -eo pid,ppid,group,user,time,vsz,rss,command'"},
		{command: "\"/bin/ps 'test' -eo pid,ppid,group,user,time,vsz,rss,command\", psutil.txt"},
		{command: "'/bin/ps \"test\" -eo pid,ppid,group,user,time,vsz,rss,command', psutil.txt", interval: 0},
		{command: "\"/bin/ps 'test' -eo pid,ppid,group,user,time,vsz,rss,command\", psutil.txt, 30, 11", interval: 30, count: 11},
		{command: "'/bin/ps \"test\" -eo pid,ppid,group,user,time,vsz,rss,command', psutil.txt, 30, 11", interval: 30, count: 11},
	}
	for _, tc := range tests {
		t.Log("TEST", tc.command)
		job := stringToJob(tc.command)
		if job == nil || !strings.Contains(job.Cmd, "/bin/ps") ||
			job.Cycle != tc.interval || job.Step != tc.count {
			t.Errorf("parse error %s", tc.command)
		}
		// t.Log(job.Cmd)
	}
}

func TestParseWorkerCommand(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigLine("STAT_CMD.Linux = 'vmstat 3 3 > _odir_/vmstat.txt'")
	schedule.ParseConfigLine("STAT_CMD.Linux = 'netstat -s', netstat.txt, 30, 10")
	if job := schedule.Collectors["Linux"].Jobs[0]; job.Cmd == "" {
		t.Error("parse worker command 1")
	}
	if job := schedule.Collectors["Linux"].Jobs[1]; job.Cycle != 30 {
		t.Error("parse worker command 2")
	}
}

func TestParseConfigFile(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigFile("../testdata/ptune/getperf.ini")
	if schedule.SiteKey != "site1" {
		t.Error("parse config 1")
	}
	colLinux := schedule.Collectors["Linux"]
	if colLinux.Jobs[0].Cmd != "/usr/bin/vmstat -a 5 61" ||
		colLinux.Jobs[1].Cmd != "/usr/bin/free -s 30 -c 12" ||
		colLinux.Jobs[2].Cmd != "/usr/bin/iostat -xk 30 12" ||
		colLinux.Jobs[3].Cmd != "/bin/cat /proc/net/dev" ||
		colLinux.Jobs[4].Cmd != "/bin/df -k -l" {
		t.Error("parse config worker job 1")
	}
}

func TestParseNGConfigFile(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigFile("../testdata/ptune_ng/getperf.ini")
	if schedule.ParseFailed == false {
		t.Error("parse wrong config")
	}
}

func TestParseLicenseFile(t *testing.T) {
	schedule := NewSchedule()
	schedule.ParseConfigFile("../testdata/ptune/network/License.txt")
	t.Log("License : ", schedule.License)
	if schedule.License.Hostname != testHost {
		t.Error("parse license")
	}
	schedule.License.Hostname = "hoge"
	schedule.ParseConfigFile("../testdata/ptune/network/License.txt")
	t.Log("License : ", schedule.License)
}
