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

	Metrics []string `toml:"metrics"`
	Servers []string `toml:"servers"`

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

# Additional metrics list of VMWare Managed Object. 
# 
# Reference : http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
# 
# example:
#
# metrics = ["parentVApp", "layout", "environmentBrowser"]

# List of monitored servers other than the local server
#
# example:
#
# servers = ["host1", "host2"]
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
