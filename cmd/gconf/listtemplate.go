package main

import (
    "fmt"
    "io"
    "os"

    // . "github.com/getperf/getperf2/common"
    "github.com/getperf/getperf2/cfg"
    "github.com/getperf/getperf2/exp"
    _ "github.com/getperf/getperf2/exp/all"
    "github.com/urfave/cli"
)

func init() {
    cmdList = append(cmdList, cli.Command{
        Name:  "ls",
        Usage: "list inventory exporter template",
        Action: func(c *cli.Context) error {
            return action(c, &listCommand{Out: os.Stdout})
        },
    })
}

type listCommand struct {
    Out io.Writer
}

func (f *listCommand) Run(c *cli.Context, r *cfg.ExportBase) error {
    for templateName := range exp.Exporters {
        fmt.Fprintf(f.Out, "\t%s\n", templateName)
    }
    return nil
}
