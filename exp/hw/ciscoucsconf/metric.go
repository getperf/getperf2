package ciscoucsconf

import "io"

type ExecType string

const (
	Cmd    = ExecType("Cmd")
	Script = ExecType("Script")
)

type Metric struct {
	Level int      `toml:"level"`
	Type  ExecType `toml:"type"`
	Id    string   `toml:"id"`
	Text  string   `toml:"text"`

	stdOut io.Writer
	stdErr io.Writer
}

type Metrics struct {
	Metrics []*Metric
}

func NewMetric(level int, execType ExecType, id string, text string) *Metric {
	metric := &Metric{
		Level: level,
		Type:  execType,
		Id:    id,
		Text:  text,
	}
	return metric
}

var metrics = []*Metric{
	NewMetric(0, "Script",
		"set_yaml",
		`set cli output yaml
		`),
}
