package primergyconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Primergy struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
	Server   string `toml:"server"`

	Metrics []*Metric `toml:"metrics"`

	url       string
	errFile   io.Writer
	vmName    string
	datastore string
	json      string
}

var sampleTemplateConfig = `
# Primergy inventory collector configuration
# Enter the information for Primergy iRMC login account and target machine
# 
# example:
#
# url = "https://192.168.10.100"  # iRMC URL or IP
# user = "test_user"
# password = "P@ssword"
# server = "w2016"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional

# Additional metrics list of Primergy REST request. 
# 
# Reference : 
# 
# example:
# 
# [[metrics]]
# 
# id = "oviewview"        # unique key
# level = 0               # command level [0-2]
# batch = ""              # profile management URL to get the batch report
# text = "/json/overview" # request URL

# The following commented out metrics are set by default

# [[metrics]]
# 
# id = "overview"
# text = "/redfish/v1/"
# 
# [[metrics]]
# 
# id = "firmware"
# text = "/redfish/v1/Systems/0/Oem/ts_fujitsu/FirmwareInventory"
# 
# [[metrics]]
# 
# id = "nic"
# text = "/redfish/v1/Systems/0/Oem/ts_fujitsu/FirmwareInventory/NIC"
# 
# [[metrics]]
# 
# id = "ntp0"
# text = "/redfish/v1/Managers/iRMC/Oem/ts_fujitsu/iRMCConfiguration/Time/NtpServers/0"
# 
# [[metrics]]
# 
# id = "ntp1"
# text = "/redfish/v1/Managers/iRMC/Oem/ts_fujitsu/iRMCConfiguration/Time/NtpServers/1"
# 
# [[metrics]]
# 
# id = "network"
# text = "/redfish/v1/Managers/iRMC/EthernetInterfaces/0"
# 
# [[metrics]]
# 
# id = "disk"
# batch = "/rest/v1/Oem/eLCM/ProfileManagement/RAIDAdapter"
# text = "/rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/HWConfigurationIrmc/Adapters/RAIDAdapter"
# 
# [[metrics]]
# 
# id = "snmp"
# batch = "/rest/v1/Oem/eLCM/ProfileManagement/NetworkServices"
# text = "/rest/v1/Oem/eLCM/ProfileManagement/get?PARAM_PATH=Server/SystemConfig/IrmcConfig/NetworkServices"
`

func (e *Primergy) Label() string {
	return "Primergy : " + e.Server
}

func (e *Primergy) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("primergyconf", func() Exporter {
		host, _ := GetHostname()
		return &Primergy{
			Server: host,
		}
	})
}
