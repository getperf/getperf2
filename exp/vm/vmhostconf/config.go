package vmwareconf

import (
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

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
# VMWare Host inventory collector configuration
# Enter the information for vCenter login account and target vm
# 
# example:
#
# url = "https://192.168.10.100/sdk"  # vCenter URL
# user = "test_user"
# password = "P@ssword"
# server = "esxi"                     # esxi host name

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
# Reference : https://code.vmware.com/apis/42/vsphere/doc/vim.HostSystem.html
# 
# [[metrics]]
# 
# id = "config"   # object key
# level = 0       # command level [0-2]
# text = "config" # If not defined, use id instead

# The following commented out metrics are set by default

# [[metrics]]
# 
# id = "summary"
# level = 0
# text = "summary"
# 
# [[metrics]]
# 
# id = "config"
# level = 0
# text = "config"

`

func (e *VMWare) Label() string {
	return "ESXi : " + e.Server
}

func (e *VMWare) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("vmhostconf", func() Exporter {
		host, _ := GetHostname()
		return &VMWare{
			Server: host,
		}
	})
}
