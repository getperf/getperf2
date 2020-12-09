package netappconf

import "io"

type ExecType string

const (
	Local  = ExecType("Local")
	Remote = ExecType("Remote")
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

var metrics = []*Metric{}
