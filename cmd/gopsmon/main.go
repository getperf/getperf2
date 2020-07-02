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
	log.Info(config)
	pid, err := config.ReadPid()
	if err != nil {
		log.Errorf("not found %s", config.PidPath)
		pid = 0
	}
	log.Info("PID:", pid)
	log.Info("OUT:", *o)

	gopscmd := `gops`
	if runtime.GOOS == "windows" {
		gopscmd = filepath.Join(config.BinDir, `gops.exe`)
	}
	log.Info("GOPS:", gopscmd)

	// gops memstats 3588
	cmdInfo1 := &agent.CommandInfo{
		CmdLine: fmt.Sprintf("%s memstats %d", gopscmd, pid),
		OutPath: filepath.Join(*o, "gops_memstats.txt"),
		Timeout: 10,
	}
	if err := cmdInfo1.ExecCommandRedirect(); err != nil {
		log.Info("ERR:", err)
	}

	// gops  3588
	cmdInfo2 := &agent.CommandInfo{
		CmdLine: fmt.Sprintf("%s %d", gopscmd, pid),
		OutPath: filepath.Join(*o, "gops.txt"),
		Timeout: 10,
	}
	if err := cmdInfo2.ExecCommandRedirect(); err != nil {
		log.Info("ERR:", err)
	}
}
