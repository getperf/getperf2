package main

import (
	"context"
	// "flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/getperf/getperf2/agent"
	_ "github.com/kardianos/minwinsvc"
)

func main() {
	log.SetFlags(0)
	subCommand := "agent"
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		subCommand = os.Args[1]
	}
	var err error
	if subCommand == "sender" {
		err = agent.RunSender(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	} else if subCommand == "admin" {
		err = agent.RunAdmin(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	} else if subCommand == "agent" {
		err = agent.Run(context.Background(), os.Args[0:], os.Stdout, os.Stderr)
	} else {
		err = fmt.Errorf("command argument must be 'sender','admin' or none %v", os.Args)
	}
	// if err != nil && err != flag.ErrHelp {
	if err != nil  {
		log.Printf("Error: %v", err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}
