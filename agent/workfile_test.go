package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func CreateTempWorkDirConfig() *Config {
	workDir, _ := ioutil.TempDir("", "")
	config := &Config{
		WorkDir:       workDir,
		WorkCommonDir: filepath.Join(workDir, "common"),
	}
	os.MkdirAll(config.WorkDir, 0777)
	os.MkdirAll(config.WorkCommonDir, 0777)
	return config
}

func CreateWorkDirTestConfig() *Config {
	workDir := "../testdata/work"
	config := &Config{
		WorkDir:       workDir,
		WorkCommonDir: filepath.Join(workDir, "common"),
	}
	return config
}

func TestGetWorkfilePath(t *testing.T) {
	config := CreateTempWorkDirConfig()
	defer os.RemoveAll(config.WorkDir)
	if path1 := config.GetWorkfilePath("test1.txt"); path1 == "" {
		t.Error("get workfile path", path1)
	}
	if path2 := config.GetWorkfilePath("_test2.txt"); path2 == "" {
		t.Error("get workfile path", path2)
	}
}

func TestWriteWorkFile(t *testing.T) {
	var tests = []struct {
		filename string
		buf      []string
		want     bool
	}{
		{"work1.txt", []string{"test1", "test2", "test3"}, true},
		{"work2/work1.txt", []string{"test1"}, false},
	}
	config := CreateTempWorkDirConfig()
	defer os.RemoveAll(config.WorkDir)

	for _, test := range tests {
		err := config.WriteWorkFile(test.filename, test.buf)
		if (err == nil && test.want == true) || (err != nil && test.want == false) {
			t.Log("write work file:", test.filename, ",expected:", test.want)
		} else {
			t.Error("write work file ", test.filename)
		}
	}
}

func TestReadWorkFile(t *testing.T) {
	config := CreateWorkDirTestConfig()
	if lines, err := config.ReadWorkFile("wk1.txt"); len(lines) != 10 {
		t.Error("read work file wk1.txt :", err, ", row : ", len(lines))
	}
	if lines, err := config.ReadWorkFile("_com1.txt"); len(lines) != 10 {
		t.Error("read work file _com1.txt :", err, ", row : ", len(lines))
	}
}

func TestReadWorkFileNumber(t *testing.T) {
	config := CreateWorkDirTestConfig()
	if i, err := config.ReadWorkFileNumber("wk_num1.txt"); i != 1234567890 {
		t.Error("read number ", err)
	}
	if i, err := config.ReadWorkFileNumber("wk1.txt"); i != 0 {
		t.Error("read illegal number ", err)
	}
}

func TestWriteWorkFileNumber(t *testing.T) {
	config := CreateTempWorkDirConfig()
	defer os.RemoveAll(config.WorkDir)
	if err := config.WriteWorkFileNumber("wk_num1.txt", 1000); err != nil {
		t.Error("write number ", err)
	}
	if i, err := config.ReadWorkFileNumber("wk_num1.txt"); i != 1000 {
		t.Error("read number ", err)
	}
}

func TestReadWorkFileHead(t *testing.T) {
	config := CreateWorkDirTestConfig()
	if lines, err := config.ReadWorkFileHead("wk1.txt", 3); err != nil && len(lines) != 3 {
		t.Error("read work file header ", err)
	}
}

func TestCheckWorkFile(t *testing.T) {
	config := CreateWorkDirTestConfig()
	if exist, err := config.CheckWorkFile("wk1.txt"); exist == false {
		t.Error("check work file ", err)
	}
	if exist, err := config.CheckWorkFile("hoge.txt"); exist == true {
		t.Error("check work file ", err)
	}
}

func TestRemoveWorkFile(t *testing.T) {
	config := CreateTempWorkDirConfig()
	defer os.RemoveAll(config.WorkDir)
	if err := config.WriteWorkFileNumber("wk2.txt", 1000); err != nil {
		t.Error("write number ", err)
	}
	if err := config.RemoveWorkFile("wk2.txt"); err != nil {
		t.Error("remove work file ", err)
	}
	if exist, err := config.CheckWorkFile("wk2.txt"); exist == true {
		t.Error("remove work file ", err)
	}
}
