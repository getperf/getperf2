package cfg

import (
	. "github.com/getperf/getperf2/common"
	"github.com/pkg/errors"
)

type CommonConfig struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Level    int    `toml:"level"`
	Timeout  int    `toml:"timeout"`
	Server   string `toml:"server"`
}

func (c *CommonConfig) Check() error {
	if c.Server == "" {
		server, err := GetHostname()
		if err != nil {
			return errors.Wrap(err, "check server parameter")
		}
		c.Server = server
	}
	return nil
}
