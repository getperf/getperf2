package solarisconf

import (
    "io"

    "github.com/getperf/getperf2/cfg"
    . "github.com/getperf/getperf2/common"
    . "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Solaris struct {
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
# Solaris server inventory collector configuration
# Enter the information for Solaris login account
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
# Solaris inventory scenarios. The text parameter using escape codes such as
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

[[metrics]]

id = "hostname_fqdn"
type = "Script"
text = '''
(
awk \'/^domain/ {print \$2}\' /etc/resolv.conf 2>/dev/null
if [ \$? != 0 ]; then
   echo 'Not Found'
fi
)
'''

[[metrics]]

id = "kernel"
text = '''
uname -X
'''

[[metrics]]

id = "cpu"
text = '''
kstat -p cpu_info
'''

[[metrics]]

id = "psrinfo"
text = '''
/usr/sbin/psrinfo
'''

[[metrics]]

id = "machineid"
text = '''
hostid
'''

[[metrics]]

id = "memory"
type = "Script"
text = '''
/usr/sbin/prtconf |grep Memory
'''

[[metrics]]

id = "swap"
text = '''
/usr/sbin/swap -s
'''

[[metrics]]

id = "network"
text = '''
/usr/sbin/ifconfig -a
'''

[[metrics]]

id = "ipadm"
text = '''
ipadm
'''

[[metrics]]

id = "net_route"
text = '''
/usr/sbin/route -v -n get default
'''

[[metrics]]

id = "ndd"
text = '''
/usr/sbin/ndd -get /dev/tcp tcp_rexmit_interval_max tcp_ip_abort_interval tcp_keepalive_interval
'''

[[metrics]]

id = "disk"
text = '''
/usr/sbin/prtpicl -v
'''

[[metrics]]

id = "metastat"
text = '''
/usr/sbin/metastat
'''

[[metrics]]

id = "filesystem"
text = '''
df -ha
'''

[[metrics]]

id = "zpool"
text = '''
/usr/sbin/zpool status
'''

[[metrics]]

id = "zpool_list"
text = '''
/usr/sbin/zpool list
'''

[[metrics]]

id = "patches"
type = "Script"
text = '''
ls /var/sadm/patch 2>/dev/null
'''

[[metrics]]

id = "solaris11_build"
text = '''
sh -c "LANG=C; /usr/bin/pkg info entire"
'''

[[metrics]]

id = "virturization"
text = '''
/usr/bin/zonename
'''

[[metrics]]

id = "packages"
text = '''
/usr/bin/pkginfo -l
'''

[[metrics]]

id = "user"
text = '''
cat /etc/passwd
'''

[[metrics]]

id = "group"
text = '''
cat /etc/group
'''

[[metrics]]

id = "service"
text = '''
sh -c "LANG=C /usr/bin/svcs -a"
'''

[[metrics]]

id = "zoneadm"
text = '''
/usr/sbin/zoneadm list -vc
'''

[[metrics]]

id = "poolstat"
text = '''
poolstat -r all
'''

[[metrics]]

id = "system_etc"
text = '''
cat /etc/system
'''

[[metrics]]

id = "resolve_conf"
text = '''
grep nameserver /etc/resolv.conf
'''

[[metrics]]

id = "coreadm"
text = '''
coreadm
'''

[[metrics]]

id = "ntp"
text = '''
egrep -e '^server' /etc/inet/ntp.conf
'''

[[metrics]]

id = "snmp_trap"
text = '''
egrep -e '^\\s*trapsink' /etc/snmp/snmpd.conf
'''
`

func (e *Solaris) Label() string {
    return "Solaris"
}

func (e *Solaris) Config() string {
    return sampleTemplateConfig
}

func init() {
    AddExporter("solarisconf", func() Exporter {
        host, _ := GetHostname()
        return &Solaris{
            Server: host,
        }
    })
}
