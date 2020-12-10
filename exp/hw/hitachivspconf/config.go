package hitachivspconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
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

	url       string
	errFile   io.Writer
	vmName    string
	datastore string
	json      string
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
# text = "/ConfigurationManager/v1/objects/storages/"
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
