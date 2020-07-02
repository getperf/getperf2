package main

import (
	"github.com/getperf/getperf2/getconfig"
)

func main() {
	if err := getconfig.Test3(); err != nil {
		panic(err)
	}
}
