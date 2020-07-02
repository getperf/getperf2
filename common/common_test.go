package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func TestCheckFile(t *testing.T) {
	var tests = []struct {
		input   string
		res     bool
		keyword string
	}{
		{"../testdata/ptune/getperf.ini", true, ""},
		{"../testdata/ptune/hogehoge", false, "not found"},
		{"../testdata/ptune", false, "is directory"},
	}
	for _, test := range tests {
		ok, err := CheckFile(test.input)
		if ok != test.res {
			t.Errorf("%q = %v,%v", test.input, test.res, test.keyword)
		}
		if test.res == false {
			if strings.Index(err.Error(), test.keyword) == -1 {
				t.Error("check file error keyword")
			}
		}
	}
}

func TestCheckDirectoryIsNull(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	if err := CheckDirectoryIsNull(home); err != nil {
		t.Error(err)
	}
	if err := CheckDirectoryIsNull("hoge"); err != nil {
		t.Log(err)
	} else {
		t.Error("check directory is null hoge")
	}
	if err := CheckDirectoryIsNull("."); err != nil {
		t.Log(err)
	} else {
		t.Error("check directory is null current dir")
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
	tmpdir, _ := ioutil.TempDir("", "example")
	defer os.Remove(tmpdir)
	var tests = []struct {
		input string
		want  bool
	}{
		{".", true},
		{tmpdir, true},
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

func TestRemoveAndCreateDir(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "example")
	defer os.Remove(tmpdir)

	path := filepath.Join(tmpdir, "hoge", "hoge")
	err := RemoveAndCreateDir(path)
	if err != nil {
		t.Error("initialize directory")
	}
	ok, _ := CheckDirectory(path)
	if !ok {
		t.Error("initialize directory")
	}
}

func TestCopyFile(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"common.go", true},
		{"hoge.go", false},
		{"soap", false},
	}
	tmpfile, err := ioutil.TempFile("", "example")
	defer os.Remove(tmpfile.Name())
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

func TestCheckDiskFree(t *testing.T) {
	if disk, err := CheckDiskFree("."); err != nil {
		t.Errorf("get disk free '.' %s", err)
	} else {
		t.Log("disk info ", disk)
	}
}
