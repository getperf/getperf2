package main

import (
	"io"
	"os"

	"github.com/getperf/getperf2/cfg"
	_ "github.com/getperf/getperf2/exp/all"
	"github.com/getperf/getperf2/gcmain"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// var env = &cfg.RunEnv{
// 	Filter:   &cfg.FilterConfig{},
// 	Retrieve: &cfg.RetrieveConfig{},
// }

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "get",
		Usage: "retrieve inventory data from getperf agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "out, o",
				Usage:       "output directory of inventory command",
				Destination: &env.Datastore,
			},
			&cli.StringFlag{
				Name:        "from-url, f",
				Usage:       "agent service url; https://agenthost:59001",
				Destination: &env.Retrieve.FromUrl,
			},
			&cli.StringFlag{
				Name:        "ca",
				Usage:       "CA certiticated file; ca.cert",
				Destination: &env.Retrieve.CaPath,
			},
			&cli.StringFlag{
				Name:        "cert",
				Usage:       "client certiticated file; client.pem",
				Destination: &env.Retrieve.ClientCertPath,
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
		},
		Action: func(c *cli.Context) error {
			return action(c, &retrieveCommand{Out: os.Stdout})
		},
	})
}

type retrieveCommand struct {
	Out io.Writer
}

func (f *retrieveCommand) Run(c *cli.Context, base *cfg.ExportBase) error {
	if err := config.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if err := env.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	controller := gcmain.NewInventoryRetriever(base, config, env)
	if controller == nil {
		return errors.New("config initialization failed")
	}
	if err := controller.Validate(); err != nil {
		return errors.Wrap(err, "run command argument check")
	}
	if err := controller.Run(); err != nil {
		return errors.Wrap(err, "run inventory executer")
	}
	return nil
}
