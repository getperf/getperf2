package eternusconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Eternus struct {
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
# Fujitsu ETERNUS storage inventory collector configuration
# Enter the information for ETERNUS login account
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
# ETERNUS inventory scenarios. The text parameter using escape codes such as
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

func (e *Eternus) Label() string {
    return "Eternus"
}

func (e *Eternus) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("eternusconf", func() Exporter {
        host, _ := GetHostname()
        return &Eternus{
            Server: host,
        }
    })
}
