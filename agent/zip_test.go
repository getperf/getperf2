package agent

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
	if err := Zip(zipPath, "../testdata", "ptune"); err != nil {
		t.Errorf("zip %s", err)
	}
	// if err := os.Remove(zipPath); err != nil {
	// 	t.Errorf("remove zip %s %s", zipPath, err)
	// }
}
