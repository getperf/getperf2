package alletraconf

import (
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Alletra struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
	Server   string `toml:"server"`

	Metrics []*Metric `toml:"metrics"`
    datastore string
}

var sampleTemplateConfig = `
# Alletra inventory collector configuration
# Enter the information for vCenter login account and target vm
# 
# example:
#
# url = "https://192.168.10.100"  # Alletra IP
# user = "test_user"
# password = "P@ssword"
# server = "w2016"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional

# Additional metrics list of Alletra Managed Object. 
# 
# Reference : 
# 
# example:
# 
# [[metrics]]
# 
# id = "arrays" # unique key
# level = 0        # command level [0-2]
# text = "arrays"

# [[metrics]]
# 
# id = "disks"
# text = "disks"
`

func (e *Alletra) Label() string {
	return "Alletra : " + e.Server
}

func (e *Alletra) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("alletraconf", func() Exporter {
		host, _ := GetHostname()
		return &Alletra{
			Server: host,
		}
	})
}
