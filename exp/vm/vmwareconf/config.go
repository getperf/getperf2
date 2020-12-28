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

var sampleTemplateConfig = `
# VMWare inventory collector configuration
# Enter the information for vCenter login account and target vm
# 
# example:
#
# url = "https://192.168.10.100/sdk"  # vCenter IP or URL
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

# Additional metrics list of VMWare vSphere Managed Object API. 
#
# Notice: 
#   The results of all added metrics are saved in the "all.json".
#
# Reference : 
#   https://code.vmware.com/apis/968/vsphere
# 
# [[metrics]]
# 
# id = "config"   # unique key
# level = 0       # command level [0-2]
# text = "config" # Managed Object Description

# The following commented out metrics are set by default

# [[metrics]]
# 
# id = "summary"
# level = 0
# text = "summary"
# 
# [[metrics]]
# 
# id = "resourceConfig"
# level = 0
# text = "resourceConfig"
# 
# [[metrics]]
# 
# id = "guestHeartbeatStatus"
# level = 0
# text = "guestHeartbeatStatus"
# 
# [[metrics]]
# 
# id = "config"
# level = 0
# text = "config"
# 
# [[metrics]]
# 
# id = "guest"
# level = 0
# text = "guest"
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
