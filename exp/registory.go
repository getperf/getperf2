package exp

import (
	"context"

	"github.com/getperf/getperf2/cfg"
)

type Exporter interface {
	// SampleConfig returns the default configuration of the Exporter
	Config() string

	// Label returns a one-sentence Label on the Input
	Label() string

	Run(ctx context.Context, env *cfg.RunEnv) error
}

type Creator func() Exporter

var Exporters = map[string]Creator{}

func AddExporter(name string, creator Creator) {
	Exporters[name] = creator
}
