package xscfconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Xscf struct {
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
# XSCF server inventory collector configuration
# Enter the information for XSCF login account
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
# XSCF inventory scenarios. The text parameter using escape codes such as
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

id = "hostname"
text = '''
uname -n
'''

`

func (e *Xscf) Label() string {
    return "Xscf"
}

func (e *Xscf) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("xscfconf", func() Exporter {
        host, _ := GetHostname()
        return &Xscf{
            Server: host,
        }
    })
}
