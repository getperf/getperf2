package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getperf/getperf2/agent"
)

func main() {
	var (
		sendFile   = flag.Bool("send", false, "send file")
		getFile    = flag.Bool("get", false, "get file")
		configFile = flag.String("config", "getperf.ini", "config file")
	)
	flag.BoolVar(sendFile, "s", false, "")
	flag.BoolVar(getFile, "g", false, "")
	flag.StringVar(configFile, "c", "", "")
	flag.Parse()

	host, err := agent.GetHostname()
	if err != nil {
		fmt.Println("get hostname failed ", err)
		os.Exit(1)
	}
	home, err := agent.GetParentAbsPath(*configFile, 1)
	if err != nil {
		fmt.Println("get home dir failed ", err)
		os.Exit(1)
	}
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("send or get file must specified")
		os.Exit(1)
	}
	dataFile := args[0]
	env := agent.NewConfigEnvBase(host, "getperfsoap", *configFile)
	config := agent.NewConfig(home, env)
	pid, err := config.ReadWorkFileNumber("_pid_getperf")
	if err != nil {
		fmt.Println("_pid_getperf not found ", err)
		os.Exit(1)
	}
	pidFile := fmt.Sprintf("_%d", pid)
	config.WorkDir = filepath.Join(home, "_wk", pidFile)
	fmt.Println("Home: ", home, ",Host: ", host)
	fmt.Println("WorkDir: ", config.WorkDir)
	if *sendFile == true {
		sendPath := config.GetWorkfilePath(dataFile)
		fmt.Println("Send : ", sendPath)
	} else if *getFile == true {
		getPath := config.GetWorkfilePath(dataFile)
		fmt.Println("Get : ", getPath)
	}
}
