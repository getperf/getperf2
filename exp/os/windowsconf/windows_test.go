package windowsconf

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/getperf/getperf2/agent"
	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/exp"
	"gopkg.in/yaml.v2"
)

var tomlpath = "../../../testdata/windowsconf.toml"
var ctx, _ = agent.MakeContext(0)

// func init() {
// 	tomlpath = filepath.Join(testNodeDir, "win2016", "windowsconf.toml")
// }

func createTestEnv() *cfg.RunEnv {
	datastore, _ := ioutil.TempDir("", "datasotre")
	env := &cfg.RunEnv{
		Level:     0,
		DryRun:    false,
		Datastore: datastore,
		LocalExec: true,
	}
	return env
}

func TestWindowsNormal(t *testing.T) {
	exp := Exporters["windowsconf"]()
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	if runtime.GOOS == "windows" {
		if err := exp.Run(ctx, env); err != nil {
			t.Error(err)
		}
		t.Log(env)
	}
}

func TestDryRunTest(t *testing.T) {
	exp := Exporters["windowsconf"]()
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	env.DryRun = true
	if runtime.GOOS == "windows" {
		if err := exp.Run(ctx, env); err != nil {
			t.Error(err)
		}
	}
}

func TestRunLevelSetTest(t *testing.T) {
	exp := Exporters["windowsconf"]()
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	env.DryRun = true
	env.Level = 1
	if runtime.GOOS == "windows" {
		if err := exp.Run(ctx, env); err != nil {
			t.Error(err)
		}
	}
}
func TestWindowsToml(t *testing.T) {
	metrics2 := Metrics{Metrics: metrics}
	d, err := yaml.Marshal(metrics2)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(d))
}

func TestWindowsConfig(t *testing.T) {
	var windows Windows
	// tomlpath := filepath.Join(testNodeDir, "win2016", "windowsconf.toml")
	_, err := toml.DecodeFile(tomlpath, &windows)
	if err != nil {
		t.Error(err)
	}
	t.Log(windows.Metrics[0])
}

func TestWindowsInventoryCode(t *testing.T) {
	var windows Windows
	env := createTestEnv()
	defer os.Remove(env.Datastore)
	// tomlpath := filepath.Join(testNodeDir, "win2016.toml")
	_, err := toml.DecodeFile(tomlpath, &windows)
	if err != nil {
		t.Error(err)
	}
	stdout := new(bytes.Buffer)
	windows.writeScript(stdout, env)
	t.Log("Result: ", stdout.String())
}
