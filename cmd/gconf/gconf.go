package main

import (
    "os"

    "github.com/getperf/getperf2/cfg"
    "github.com/getperf/getperf2/gcmain"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "github.com/urfave/cli"
)

type subCmd interface {
    Run(*cli.Context, *cfg.ExportBase) error
}

var cmdList = []cli.Command{}

var base = &cfg.ExportBase{}

func main() {
    app := cli.NewApp()
    app.Name = "gconf"
    app.Usage = "Getconfig inventory exporter"
    app.Version = gcmain.Version

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:        "template, t",
            Usage:       "inventory collector template name",
            Destination: &base.Template,
        },
        cli.StringFlag{
            Name:        "config, c",
            Usage:       "config path of template",
            Destination: &base.ConfigPath,
        },
    }

    app.Commands = cmdList
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func action(c *cli.Context, sc subCmd) error {
    var err error
    if base.Check() != nil {
        return errors.Wrap(err, "check config")
    }
    // base.Home = GetBaseDir()
    // if configPath := base.ConfigPath; configPath == "" {
    //     configName := fmt.Sprintf("%s.toml", base.Template)
    //     base.ConfigPath = filepath.Join(base.Home, configName)
    // } else {
    //     base.Home, err = GetParentAbsPath(configPath, 1)
    //     if err != nil {
    //         return errors.Wrap(err, "check config")
    //     }
    // }
    return sc.Run(c, base)
}
