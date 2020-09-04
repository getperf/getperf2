package vmwareconf

import (
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Metric struct {
	Id    string `toml:"id"`
	Level int    `toml:"level"`
	Text  string `toml:"text"`
}

func (metric *Metric) getObjectId() string {
	if metric.Text == "" {
		return metric.Id
	} else {
		return metric.Text
	}
}

type VMWare struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
	Server   string `toml:"server"`

	Servers []string  `toml:"servers"`
	Metrics []*Metric `toml:"metrics"`

	vmName    string
	datastore string
	json      string
}

type Command struct {
	Id string `toml:"id"`
}

var sampleTemplateConfig = `
# VMWare inventory collector configuration
# Enter the information for vCenter login account and target vm
# 
# example:
#
# url = "https://192.168.10.100/sdk"  # vCenter URL
# user = "test_user"
# password = "P@ssword"
# server = "w2016"                    # vm name

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

# Additional metrics list of VMWare Managed Object. 
#
# Notice: 
#   The results of all added metrics are saved in the "all.json".
#
# Reference : http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
# 
# [[metrics]]
# 
# id = "config"   # object key
# level = 0       # command level [0-2]
# text = "config" # If not defined, use id instead

[[metrics]]

id = "config"

[[metrics]]

id = "summary"

[[metrics]]

id = "capability"

[[metrics]]

id = "datastore"

[[metrics]]

id = "environmentBrowser"

[[metrics]]

id = "guest"

[[metrics]]

id = "guestHeartbeatStatus"

[[metrics]]

id = "layoutEx"

[[metrics]]

id = "network"

[[metrics]]

id = "resourceConfig"

[[metrics]]

id = "resourcePool"

[[metrics]]

id = "snapshot"

[[metrics]]

id = "storage"

[[metrics]]

id = "runtime"

`

func (e *VMWare) Label() string {
	return "VMWare : " + e.Server
}

func (e *VMWare) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("vmwareconf", func() Exporter {
		host, _ := GetHostname()
		return &VMWare{
			Server: host,
		}
	})
}
