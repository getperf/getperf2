package gcmain

import (
	"io/ioutil"
	_ "os"
	"path/filepath"
	"testing"
)

func TestZip(t *testing.T) {
	logDir, _ := ioutil.TempDir("", "log")
	// defer os.RemoveAll(logDir)

	zipPath := filepath.Join(logDir, "test.zip")
	t.Log("Zip:", zipPath)
	if err := zipForAgent(zipPath, "../testdata", "hos1/WindowsConf/20210129/150000"); err != nil {
		t.Errorf("zip %s", err)
	}
}
