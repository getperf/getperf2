package hitachivspconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
	"github.com/go-resty/resty/v2"
)

const Version = "0.1.4"

type HitachiVSP struct {
	Url      string    `toml:"url"`
	User     string    `toml:"user"`
	Password string    `toml:"password"`
	Insecure bool      `toml:"insecure"`
	Server   string    `toml:"server"`
	Servers  []string  `toml:"servers"`
	Metrics  []*Metric `toml:"metrics"`

	url             string
	errFile         io.Writer
	vmName          string
	datastore       string
	json            string
	client          *resty.Client
	storageDeviceId string
	token           string
	sessionId       string
}

var sampleTemplateConfig = `
# HitachiVSP inventory collector configuration
# Enter the information for HitachiVSP REST login account and target machine
# 
# example:
#
# url = "https://192.168.10.100"  # REST URL or IP
# user = "test_user"
# password = "P@ssword"
# server = "w2016"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional

# List of monitored servers other than the local server
#
# example:
#
# servers = ["storage1", "storage2"]

# Additional metrics list of HitachiVSP REST request. 
# 
# Reference : 
# 
# http://itdoc.hitachi.co.jp/manuals/3021/3021901610/INDEX.HTM
# 
# example:
# 
# [[metrics]]
# 
# id = "oviewview"        # unique key
# level = 0               # command level [0-2]
# text = "/json/overview" # request URL. Replace the string '{id}' with the storage device ID

# The following commented out metrics are set by default

# [[metrics]]
# 
# id = "storage"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}"
#
# [[metrics]]
# 
# id = "host-groups"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/host-groups"
#
# [[metrics]]
# 
# id = "ports"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/ports"
#
# [[metrics]]
# 
# id = "ports"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/ports"
#
# [[metrics]]
# 
# id = "parity-groups"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/parity-groups"
#
# [[metrics]]
# 
# id = "ldevs"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/ldevs?headLdevId=0&count=100"
#
# [[metrics]]
# 
# id = "users"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/users"
#
# [[metrics]]
# 
# id = "ambient"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/components/instance"
#
# [[metrics]]
# 
# id = "snmp"
# level = 0
# text = "/ConfigurationManager/v1/objects/storages/{id}/snmp-settings/instance"
`

func (e *HitachiVSP) Label() string {
	return "HitachiVSP : " + e.Server
}

func (e *HitachiVSP) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("hitachivspconf", func() Exporter {
		host, _ := GetHostname()
		return &HitachiVSP{
			Server: host,
		}
	})
}
