package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewDatastoreBase(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	ds := NewDatastoreBase(home, "hoge", "HW", time.Now())
	absDir := strings.TrimPrefix(ds.AbsDir(), home)
	if len(absDir) == 0 {
		t.Error("new out log base is nil")
	}
	if TrimPathSeparator(absDir) != ds.RelDir() {
		t.Error("new out log base")
	}
}

func TestNewDatastore(t *testing.T) {
	home, _ := ioutil.TempDir("", "ptune")
	defer os.RemoveAll(home)
	config := NewConfig(home, NewConfigEnv())
	config.InitAgent()
	config.ParseConfigLine("STAT_INTERVAL.HW = 300")
	config.ParseConfigLine("STAT_CMD.HW = 'netstat -s', netstat.txt")

	date := "2020-01-20 11:00:00"
	now, _ := time.Parse("2006-01-02 15:04:05", date)

	ds, err := config.NewDatastore("HW", now)
	if err != nil {
		t.Error("new out log", err)
	}
	if ds.ZipFile("hoge") != "arc_hoge__HW_20200120_110000.zip" {
		t.Error("new out log 2")
	}
}

func TestOldestZip(t *testing.T) {
	ds := NewDatastoreBase("../testdata/ptune/backup",
		"centos80", "HW", time.Now())
	if zip, _ := ds.OldestZip(""); zip != "arc_centos80__HW_20200528_132000.zip" {
		t.Error("test oldest zip")
	}
	if zip, _ := ds.OldestZip("20200528_132000"); zip != "arc_centos80__HW_20200528_132500.zip" {
		t.Error("test oldest zip2")
	}
	zip, err := ds.OldestZip("20200529_132000")
	if zip != "" || err == nil {
		t.Error("test oldest zip3")
	}
	t.Log(err)
}

func TestLatestZip(t *testing.T) {
	ds := NewDatastoreBase("../testdata/ptune/backup",
		"centos80", "HW", time.Now())
	zip, err := ds.LatestZip()
	if zip != "arc_centos80__HW_20200528_133500.zip" {
		t.Error("test latest zip")
	}
	if err != nil {
		t.Error("test latest zip :", err)
	}
}

func TestDatastoreKeys(t *testing.T) {
	ds := NewDatastoreBase("../testdata/ptune/backup",
		"centos80", "HW", time.Now())
	datastoreSetes, err := ds.GetDatastoreSets()
	if err != nil {
		t.Error("test datastore zip key :", err)
	}
	jsonBytes, err := json.Marshal(datastoreSetes)
	if err != nil {
		t.Error("JSON Marshal error:", err)
	}
	t.Logf("|%s|", string(jsonBytes))
}
