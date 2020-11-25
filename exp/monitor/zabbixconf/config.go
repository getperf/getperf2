package zabbixconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
	"github.com/rday/zabbix"
)

const Version = "0.1.4"

type Zabbix struct {
	Url      string   `toml:"url"`
	User     string   `toml:"user"`
	Password string   `toml:"password"`
	Insecure bool     `toml:"insecure"`
	Server   string   `toml:"server"`
	Servers  []string `toml:"servers"`

	datastore string
	errFile   io.Writer
	session   *zabbix.API
	HostIds   map[string]string
}

var sampleTemplateConfig = `
# Zabbix monitor inventory collector configuration
# Enter the information for Zabbix login account and target server
# 
# example:
#
# url = "https://192.168.10.100/zabbix"  # Zabbix URL
# user = "Admin"
# password = "P@ssword"
# server = "centos80"                    # monitoring host name

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
server = "{{ .Server }}"

# The following parameters are optional

# List of monitored servers other than the local server
#
# example:
#
# servers = ["host1", "host2"]
`

func (e *Zabbix) Label() string {
	return "Zabbix : " + e.Server
}

func (e *Zabbix) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("zabbixconf", func() Exporter {
		host, _ := GetHostname()
		return &Zabbix{
			Server: host,
		}
	})
}
