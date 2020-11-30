package ciscoucsconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type CiscoUCS struct {
    Server     string `toml:"server"`
    Url        string `toml:"url"`
    User       string `toml:"user"`
    Password   string `toml:"password"`
    SshKeyPath string `toml:"ssh_key"`
    Insecure   bool   `toml:"insecure"`

    Metrics []*Metric `toml:"metrics"`

    Env       *cfg.RunEnv
    errFile   io.Writer
    datastore string
}

var sampleTemplateConfig = `
# CiscoUCS server inventory collector configuration
# Enter the information for CiscoUCS login account
# 
# example:
#
# url = "192.168.10.100"
# user = "test_user"
# password = "P@ssword"
# server = "sol10"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional


# Describe the additional command list. Added to the default command list for
# CiscoUCS inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
# 
# [[metrics]]
# 
# id = "oracle_module"   # unique key
# level = 0    # command level [0:Default,1,2]
# type = "Cmd" # "Cmd":single command(Defalut), "Script":multi line commands
# text = '''
# ls /home/oracle/"
# '''

[[metrics]]

id = "bios"
type = "Script"
text = '''
top
show bios detail
'''

[[metrics]]

id = "chassis"
type = "Script"
text = '''
top
show chassis detail
'''

[[metrics]]

id = "cimc"
type = "Script"
text = '''
top
show cimc detail
'''

[[metrics]]

id = "cpu"
type = "Script"
text = '''
top
scope chassis
show cpu detail
'''

[[metrics]]

id = "memory"
type = "Script"
text = '''
top
scope chassis
show dimm-summary
'''

[[metrics]]

id = "hdd"
type = "Script"
text = '''
top
scope chassis
show  hdd-pid detail
'''

[[metrics]]

id = "storageadapter"
type = "Script"
text = '''
top
scope chassis
show  storageadapter detail
'''

[[metrics]]

id = "physical_drive"
type = "Script"
text = '''
top
scope chassis
scope storageadapter MRAID
show  physical-drive detail
'''

[[metrics]]

id = "virtual_drive"
type = "Script"
text = '''
top
scope chassis
scope storageadapter MRAID
show  virtual-drive detail
'''

[[metrics]]

id = "network"
type = "Script"
text = '''
top
show cimc/network detail
'''

[[metrics]]

id = "snmp"
type = "Script"
text = '''
top
show snmp detail
'''

[[metrics]]

id = "snmp_trap"
type = "Script"
text = '''
top
show snmp/trap-destinations detail
'''

[[metrics]]

id = "ntp"
type = "Script"
text = '''
top
show /cimc/network/ntp detail
'''
`

func (e *CiscoUCS) Label() string {
    return "CiscoUCS"
}

func (e *CiscoUCS) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("ciscoucsconf", func() Exporter {
        host, _ := GetHostname()
        return &CiscoUCS{
            Server: host,
        }
    })
}
