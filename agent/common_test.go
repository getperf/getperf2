package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetHostnameAlias(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"hostA", "hosta"},
		{"hostb.getperf.com", "hostb"},
		{"", ""},
	}
	for _, test := range tests {
		get, err := GetHostnameAlias(test.input)
		t.Log(test.input, " -> ", get)
		if err != nil || get != test.want {
			t.Errorf("GetHostname(%q) = %v", test.input, test.want)
		}
	}
}

func TestGetHostname(t *testing.T) {
	get, err := GetHostname()
	t.Log(get)
	if err != nil || len(get) < 1 {
		t.Errorf("GetHostname() = %v", get)
	}
}

func TestGetParentPath(t *testing.T) {
	var tests = []struct {
		input string
		level int
		want  string
	}{
		{"/foo/bar/baz", 1, filepath.Clean("/foo/bar")},
		{"/foo/bar/baz", 2, filepath.Clean("/foo")},
		{"/foo/bar/baz", 3, filepath.Clean("/")},
		{"/foo/bar/baz", 4, filepath.Clean("/")},
	}
	for _, test := range tests {
		get := GetParentPath(test.input, test.level)
		t.Log(test.input, " -> ", get)
		if get != test.want {
			t.Errorf("GetParentPath(%q, %d) = %v", test.input, test.level, test.want)
		}
	}
}

func TestGetParentAbsPath(t *testing.T) {
	get, err := GetParentAbsPath(".", 1)
	t.Log(get)
	if err != nil || len(get) < 1 {
		t.Errorf("GetHostname() = %v", get)
	}
}

func TestCheckDirectory(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{".", true},
		{"/foo/bar/baz", false},
		{"common.go", false},
	}
	for _, test := range tests {
		get, err := CheckDirectory(test.input)
		t.Log(test.input, " -> ", get, err)
		if get != test.want {
			t.Errorf("CheckDirectory(%q) = %v", test.input, test.want)
		}
	}
}

func TestCopyFile(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"common.go", true},
		{"config.go", true},
		{"hoge.go", false},
		{"soap", false},
	}
	tmpfile, err := ioutil.TempFile("", "example")
	defer os.Remove(tmpfile.Name()) // clean up
	if err != nil {
		t.Errorf("create tempfile")
	}
	t.Log("careate target", tmpfile.Name())
	for _, test := range tests {
		err := CopyFile(test.input, tmpfile.Name())
		if (err == nil && test.want == true) || (err != nil && test.want == false) {
			t.Log("copy file ", test.input, test.want)
		} else {
			t.Error("copy file ", test.input)
		}
	}
}

func TestGetTimeString(t *testing.T) {
	now := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	if str := GetTimeString(DEFAULT, now); str != "2018/01/01 00:00:00" {
		t.Errorf("get time default")
	}
	if str := GetTimeString(YYYYMMDD, now); str != "20180101" {
		t.Errorf("get time YYYYMMDD")
	}
	if str := GetTimeString(DIR, now); str != filepath.Join("20180101", "000000") {
		t.Errorf("get time DIR")
	}
}

func TestGetCurrentTime(t *testing.T) {
	if str := GetCurrentTime(3600, DEFAULT); len(str) != 19 {
		t.Errorf("get time default %s %d", str, len(str))
	}
}

func TestCheckPathInHome(t *testing.T) {
	home, _ := filepath.Abs(".")
	config := &Config{Home: home}
	if isInclude, err := config.CheckPathInHome("."); isInclude == false {
		t.Errorf("check path in home %v : %s", isInclude, err)
	}
}

func TestCheckDiskFree(t *testing.T) {
	if disk, err := CheckDiskFree("."); err != nil {
		t.Errorf("get disk free '.' %s", err)
	} else {
		t.Log("disk info ", disk)
	}
}
