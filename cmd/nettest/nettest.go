package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/getperf/getperf2/agent"
)

func main() {
	log.SetFlags(0)
	err := agent.RunNettest(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}
