package netappconf

import "io"

type ExecType string

const (
	Cmd    = ExecType("Cmd")
	Script = ExecType("Script")
)

type Metric struct {
	Level  int      `toml:"level"`
	Type   ExecType `toml:"type"`
	Remote bool     `toml:"remote"`
	Id     string   `toml:"id"`
	Text   string   `toml:"text"`

	stdOut io.Writer
	stdErr io.Writer
}

type Metrics struct {
	Metrics []*Metric
}

func NewMetric(level int, execType ExecType, remote bool, id string, text string) *Metric {
	metric := &Metric{
		Level:  level,
		Type:   execType,
		Remote: remote,
		Id:     id,
		Text:   text,
	}
	return metric
}

var metrics = []*Metric{
	NewMetric(0, "Cmd", true, "subsystem_health", "system health subsystem show -node {host}"),
	NewMetric(0, "Cmd", true, "storage_failover", "storage failover show -node {host}"),
	NewMetric(0, "Cmd", true, "memory", "system controller memory dimm show -node {host}"),
	NewMetric(0, "Cmd", true, "license", "system license show -owner {host}"),
	NewMetric(0, "Cmd", true, "processor", "system controller show -node {host}"),
	NewMetric(0, "Cmd", true, "volume", "volume show -nodes {host}"),
	NewMetric(0, "Cmd", true, "aggregate_status", "aggr show -owner-name {host}"),
	NewMetric(0, "Cmd", true, "sysconfig", "run {host} sysconfig -a"),
	NewMetric(0, "Cmd", true, "sysconfig_raid", "run {host} sysconfig -r"),
	NewMetric(0, "Cmd", true, "network_interface", "network interface show -curr-node {host}"),
	NewMetric(0, "Cmd", false, "version", "version"),
	NewMetric(0, "Cmd", false, "vserver", "vserver show"),
	NewMetric(0, "Cmd", false, "snmp", "system snmp show"),
	NewMetric(0, "Cmd", false, "ntp", "cluster time-service ntp server show"),
	NewMetric(0, "Cmd", false, "df", "df"),
}
