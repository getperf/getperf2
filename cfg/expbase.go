package cfg

import (
	"fmt"
	"path/filepath"

	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ExportBase struct {
	Home       string
	Template   string
	ConfigPath string
}

func (b *ExportBase) Check() error {
	if b.Home == "" {
		b.Home = GetBaseDir()
	}
	if configPath := b.ConfigPath; configPath == "" {
		configName := fmt.Sprintf("%s.toml", b.Template)
		b.ConfigPath = filepath.Join(b.Home, configName)
	} else {
		home, err := GetParentAbsPath(configPath, 1)
		if err != nil {
			return errors.Wrap(err, "check config")
		}
		b.Home = home
	}
	log.Info("set base home : ", b.Home)
	return nil
}
