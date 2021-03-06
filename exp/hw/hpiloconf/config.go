package hpiloconf

import (
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type HPiLO struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
	Server   string `toml:"server"`

	Metrics []*Metric `toml:"metrics"`

	vmName    string
	datastore string
	json      string
}

var sampleTemplateConfig = `
# HPiLO inventory collector configuration
# Enter the information for vCenter login account and target vm
# 
# example:
#
# url = "https://192.168.10.100"  # iLO URL or IP
# user = "test_user"
# password = "P@ssword"
# server = "w2016"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional

# Additional metrics list of HPiLO Managed Object. 
# 
# Reference : 
# 
# example:
# 
# [[metrics]]
# 
# id = "oviewview" # unique key
# level = 0        # command level [0-2]
# text = "/json/overview"

# [[metrics]]
# 
# id = "overview"
# text = "/json/overview"
# 
# [[metrics]]
# 
# id = "license"
# text = "/redfish/v1/Managers/1/LicenseService/1/"
# 
# [[metrics]]
# 
# id = "proc_info"
# text = "/json/proc_info"
# 
# [[metrics]]
# 
# id = "mem_info"
# text = "/json/mem_info"
# 
# [[metrics]]
# 
# id = "network"
# text = "/redfish/v1/Managers/1/EthernetInterfaces/1"
# 
# [[metrics]]
# 
# id = "health_drives"
# text = "/json/health_drives"
# 
# [[metrics]]
# 
# id = "health_phy_drives"
# text = "/json/health_phy_drives"
# 
# [[metrics]]
# 
# id = "snmp"
# text = "/redfish/v1/Managers/1/snmpservice/snmpalertdestinations/1"
# 
# [[metrics]]
# 
# id = "power_regulator"
# text = "/json/power_regulator"
# 
# [[metrics]]
# 
# id = "power_summary"
# text = "/json/power_summary"
`

func (e *HPiLO) Label() string {
	return "HPiLO : " + e.Server
}

func (e *HPiLO) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("hpiloconf", func() Exporter {
		host, _ := GetHostname()
		return &HPiLO{
			Server: host,
		}
	})
}
