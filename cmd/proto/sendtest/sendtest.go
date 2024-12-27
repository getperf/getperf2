package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/getperf/getperf2/agent"
	log "github.com/sirupsen/logrus"
)

const cmdName = "getperf"

func initConfig(configFile string) (*agent.Config, error) {
	hostName, err := agent.GetHostname()
	log.Info("Host:", hostName)
	if err != nil {
		return nil, fmt.Errorf("gethost %s", err)
	}
	if configFile == "" {
		home, err := agent.GetParentAbsPath(os.Args[0], 2)
		if err != nil {
			return nil, fmt.Errorf("get program path %s", err)
		}
		configFile = filepath.Join(home, "getperf.ini")
	}
	configEnv := agent.NewConfigEnvBase(hostName, cmdName, configFile)
	home, err := agent.GetParentAbsPath(configFile, 1)
	if err != nil {
		return nil, fmt.Errorf("get parent path %s", err)
	}
	config := agent.NewConfig(home, configEnv)
	return config, nil
}

func main() {
	var (
		o = flag.String("o", "", "output directory")
		c = flag.String("c", "", "config file")
	)
	flag.Parse()
	config, err := initConfig(*c)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("OUT:", *o)

	getperf2 := filepath.Join(config.BinDir, `getperf2`)
	if runtime.GOOS == "windows" {
		getperf2 = filepath.Join(config.BinDir, `getperfsoap.exe`)
	}
	log.Info("cmd:", getperf2)

	// gops memstats 3588
	cmdInfo1 := &agent.CommandInfo{
		// CmdLine: fmt.Sprintf("%s  --send aaa", getperf2),
		CmdLine: fmt.Sprintf("%s  --send arc_aaa__Linux_20230808_1100.zip", getperf2),
		OutPath: filepath.Join(*o, "sendtest.txt"),
		Timeout: 10,
	}
	if err := cmdInfo1.ExecCommandRedirect(); err != nil {
		log.Info("ERR:", err)
	}
}
