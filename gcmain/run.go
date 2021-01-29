package gcmain

import (
	"context"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/getperf/getperf2/cfg"
	"github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/exp"
	_ "github.com/getperf/getperf2/exp/all"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type InventoryExecuter struct {
	Base   *cfg.ExportBase
	Common *cfg.CommonConfig
	Env    *cfg.RunEnv
}

func NewInventoryExecuter(base *cfg.ExportBase, common *cfg.CommonConfig, env *cfg.RunEnv) *InventoryExecuter {
	if base == nil || common == nil || env == nil {
		return nil
	}
	return &InventoryExecuter{
		Base:   base,
		Common: common,
		Env:    env,
	}
}

func (c *InventoryExecuter) Validate() error {
	if err := c.Common.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if err := c.Env.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if c.Base.Template == "" {
		return errors.New("--teplate must specifiend")
	}
	return nil
}

func MakeContext(timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		duration := time.Duration(timeout) * time.Second
		return context.WithTimeout(context.Background(), duration)
	} else {
		return context.WithCancel(context.Background())
	}
}

func (c *InventoryExecuter) Run() error {
	ctx, cancel := MakeContext(c.Env.Timeout)
	defer cancel()
	if err := common.SetLogLevel(c.Env.LogLevel); err != nil {
		return fmt.Errorf("set log level : %s.", err)
	}
	log.Info("set inventory datastore : ", c.Env.Datastore)
	templateName := c.Base.Template
	scenario, ok := exp.Exporters[templateName]
	if !ok {
		return fmt.Errorf("unkown template : %s.", templateName)
	}
	log.Debugf("run (label=%s, template=%s)\n", scenario().Label(), templateName)
	server := scenario()
	_, err := toml.DecodeFile(c.Base.ConfigPath, server)
	if err != nil {
		return errors.Wrap(err, "read server config file")
	}
	if err := server.Run(ctx, c.Env); err != nil {
		return err
	}
	if c.Env.Send {
		sender := NewSender(scenario().Label(), c.Env)
		if err := sender.Run(); err != nil {
			return errors.Wrap(err, "end processing")
		}
	}
	return nil
}
