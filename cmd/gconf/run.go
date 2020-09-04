package main

import (
	"io"
	"os"

	// . "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/cfg"
	_ "github.com/getperf/getperf2/exp/all"
	"github.com/getperf/getperf2/gcmain"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var env = &cfg.RunEnv{
	Filter:   &cfg.FilterConfig{},
	Retrieve: &cfg.RetrieveConfig{},
}

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "run",
		Usage: "run inventory exporter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "out, o",
				Usage:       "output directory of inventory command",
				Destination: &env.Datastore,
			},
			&cli.BoolFlag{
				Name:        "dryrun, d",
				Usage:       "use dry run mode",
				Destination: &env.DryRun,
			},
			cli.IntFlag{
				Name:        "level, l",
				Usage:       "run level",
				Destination: &env.Level,
			},
			cli.IntFlag{
				Name:        "timeout",
				Usage:       "command timeout sec",
				Destination: &env.Timeout,
			},
			cli.IntFlag{
				Name:        "log-level",
				Usage:       "log level[0-7]",
				Value:       6,
				Destination: &env.LogLevel,
			},
			&cli.StringFlag{
				Name:        "keyword-server",
				Usage:       "filtering keyword of target servers",
				Destination: &env.Filter.KeywordServer,
			},
			&cli.StringFlag{
				Name:        "keyword-testitem",
				Usage:       "filtering keyword of test itemss",
				Destination: &env.Filter.KeywordTestItem,
			},
		},
		Action: func(c *cli.Context) error {
			return action(c, &runCommand{Out: os.Stdout})
		},
	})
}

type runCommand struct {
	Out io.Writer
}

func (f *runCommand) Run(c *cli.Context, base *cfg.ExportBase) error {
	if err := config.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if err := env.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	controller := gcmain.NewInventoryExecuter(base, config, env)
	if controller == nil {
		return errors.New("config initialization failed")
	}
	if err := controller.Validate(); err != nil {
		return errors.Wrap(err, "run command argument check")
	}
	if err := controller.Run(); err != nil {
		return errors.Wrap(err, "run inventory executer")
	}

	// log.Printf("run config : %v\n", config)
	// log.Printf("run base : %v\n", base)
	// log.Printf("run env : %v, filter: %v\n", env, env.Filter)
	// if templateName == "" {
	// 	return errors.New("--template required. Usage: gconf --template {xxxconf} init ...")
	// }
	// exp, ok := exporter.Exporters[templateName]
	// if !ok {
	// 	return fmt.Errorf("unkown template : %s.", templateName)
	// }
	// baseDir := GetBaseDir()
	// configPath := filepath.Join(baseDir, fmt.Sprintf("%s.toml", templateName))
	// if c.String("config") != "" {
	// 	configPath = c.String("config")
	// }
	// fmt.Fprintf(f.Out, "CONFIGPATH :%s\n", configPath)
	// fmt.Fprintf(f.Out, "run (label=%s, template=%s)\n", exp().Label(), templateName)
	// server := exp()
	// _, err := toml.DecodeFile(configPath, server)
	// if err != nil {
	// 	return errors.Wrap(err, "making server")
	// }
	// fmt.Fprintf(f.Out, "run (label=%s, template=%s)\n", server.Label(), templateName)
	// fmt.Fprintf(f.Out, "env=%v\n", env)
	// return server.Run(&env)
	return nil
}

// func InitConfig() error {
//     baseDir := GetBaseDir()
//     log.Info("set base directory ", baseDir)
//     configPath := filepath.Join(baseDir, "windowsconf.toml")
//     return InitConfigFile(configPath)
// }

// func InitConfigFile(configPath string) error {
//     //     if ok, _ := CheckFile(configPath); ok {
//     //         log.Warnf("'%s' exist, Backup to '%s'", configPath, configPath+"_bak")
//     //         if err := CopyFile(configPath, configPath+"_bak"); err != nil {
//     //             return errors.Wrap(err, "backup config file")
//     //         }
//     //     }
//     //     contents := []byte(Config())
//     //     if err := ioutil.WriteFile(configPath, contents, 0664); err != nil {
//     //         return errors.Wrap(err, "write config")
//     //     }
//     return nil
// }
