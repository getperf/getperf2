package gcmain

import (
	"fmt"
	"text/template"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	"github.com/getperf/getperf2/exp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ConfigInitializer struct {
	Base   *cfg.ExportBase
	Common *cfg.CommonConfig
}

func NewConfigInitializer(base *cfg.ExportBase, common *cfg.CommonConfig) *ConfigInitializer {
	if base == nil || common == nil {
		return nil
	}
	return &ConfigInitializer{
		Base:   base,
		Common: common,
	}
}

func (c *ConfigInitializer) Validate() error {
	if err := c.Common.Check(); err != nil {
		return errors.Wrap(err, "check config")
	}
	if c.Base.Template == "" {
		return errors.New("--teplate must specifiend")
	}
	return nil
}

func (c *ConfigInitializer) CreateConfig() error {
	templateName := c.Base.Template
	exp, ok := exp.Exporters[templateName]
	if !ok {
		return fmt.Errorf("unkown template : %s.", templateName)
	}
	configText := exp().Config()
	tpl, err := template.New("config").Parse(configText)
	if err != nil {
		return errors.Wrap(err, "initialize config template")
	}
	configPath := c.Base.ConfigPath
	configFile, err := CreateAndOpenFile(configPath)
	if err != nil {
		return errors.Wrap(err, "create config file")
	}
	defer configFile.Close()

	err = tpl.Execute(configFile, c.Common)
	if err != nil {
		return errors.Wrap(err, "make config file")
	}
	log.Info("config created : ", configPath)

	return nil
}
