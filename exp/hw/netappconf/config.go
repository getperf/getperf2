package netappconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type NetAPP struct {
    Url        string    `toml:"url"`
    User       string    `toml:"user"`
    Password   string    `toml:"password"`
    SshKeyPath string    `toml:"ssh_key"`
    Insecure   bool      `toml:"insecure"`
    Server     string    `toml:"server"`
    Servers    []string  `toml:"servers"`
    Metrics    []*Metric `toml:"metrics"`

    Env       *cfg.RunEnv
    errFile   io.Writer
    datastore string
}

var sampleTemplateConfig = `
# NetAPP storage inventory collector configuration
# Enter the information for NetAPP CLI SSH login account
# 
# example:
#
# url = "192.168.10.100"
# user = "test_user"
# password = "P@ssword"
# server = "netapp"

url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
server = "{{ .Server }}"

# The following parameters are optional


# Describe the additional command list. Added to the default command list for
# NetAPP inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
#
# [[metrics]]
# 
# id = "sysconfig"       # unique key
# level = 0              # command level [0:Default,1,2]
# remote = true          # Set to false for clusters, true for remote server
# type = "Cmd"           # "Cmd" or "Shell"
# text = '''             # Command. Replace the string '{host}' with the remote server
# run {host} sysconfig -a
# '''

# The following commented out metrics are set by default

# [[metrics]]
# 
# id = "subsystem_health"
# remote = true
# text = '''
# system health subsystem show -node {host}
# '''
# 
# [[metrics]]
# 
# id = "storage_failover"
# remote = true
# text = '''
# storage failover show -node {host}
# '''
# 
# [[metrics]]
# 
# id = "memory"
# remote = true
# text = '''
# system controller memory dimm show -node {host}
# '''
# 
# [[metrics]]
# 
# id = "license"
# remote = true
# text = '''
# system license show -owner {host}
# '''
# 
# [[metrics]]
# 
# id = "processor"
# remote = true
# text = '''
# system controller show -node {host}
# '''
# 
# [[metrics]]
# 
# id = "volume"
# remote = true
# text = '''
# volume show -nodes {host}
# '''
# 
# [[metrics]]
# 
# id = "aggregate_status"
# remote = true
# text = '''
# aggr show -owner-name {host}
# '''
# 
# [[metrics]]
# 
# id = "sysconfig"
# remote = true
# text = '''
# run {host} sysconfig -a
# '''
# 
# [[metrics]]
# 
# id = "sysconfig_raid"
# remote = true
# text = '''
# run {host} sysconfig -r
# '''
# 
# [[metrics]]
# 
# id = "snmp"
# remote = false
# text = '''
# system snmp show
# '''
# 
# [[metrics]]
# 
# id = "ntp"
# remote = false
# text = '''
# cluster time-service ntp server show
# '''
# 
# [[metrics]]
# 
# id = "version"
# remote = false
# text = '''
# version
# '''
#
# [[metrics]]
# 
# id = "network_interface"
# remote = true
# text = '''
# network interface show -curr-node {host}
# '''
# 
# [[metrics]]
# 
# id = "vserver"
# remote = false
# text = '''
# vserver show
# '''
#
# [[metrics]]
# 
# id = "df"
# remote = false
# text = '''
# df
# '''
`

func (e *NetAPP) Label() string {
    return "NetAPP"
}

func (e *NetAPP) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("netappconf", func() Exporter {
        host, _ := GetHostname()
        return &NetAPP{
            Server: host,
        }
    })
}
