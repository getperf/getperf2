package linuxconf

import (
	"io"

	"github.com/getperf/getperf2/cfg"
	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Linux struct {
	Server    string    `toml:"server"`
	LocalExec bool      `toml:"local_exec"`
	Servers   []*Server `toml:"servers"`
	Metrics   []*Metric `toml:"metrics"`

	Env          *cfg.RunEnv
	errFile      io.Writer
	remoteServer string
	datastore    string
}

type Server struct {
	Server     string `toml:"server"`
	Url        string `toml:"url"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	SshKeyPath string `toml:"ssh_key"`
	Insecure   bool   `toml:"insecure"`
}

var sampleTemplateConfig = `
# Linux server inventory collector configuration
# When collecting the inventory of Linux platform, execute it locally.
# Therefore, no template setting is required

server = "{{ .Server }}"
local_exec = true

# The following parameters are optional

# Enter the information for login account of remote linux server
# 
# example:
# 
# [[servers]]
# 
# server = "centos8"    # server name
# url = "192.168.10.1"  # server address, example: 192.168.0.1 , 192.168.0.1:22
# user = "test_user"
# password = "P@ssword"
# ssh_key = ""          # ssh private key path, ignore if it not set

# Describe the additional command list. Added to the default command list for
# Linux inventory scenarios. The text parameter using escape codes such as
# '\"', '\\', See these example,
# 
# example:
# 
# [[metrics]]
# 
# id = "oracle_module"   # unique key
# level = 0      # command level [0-2]
# text = "ls /home/oracle/"

[[metrics]]

id = "hostname"
level = 0
text = '''
hostname -s
'''

[[metrics]]

id = "hostname_fqdn"
level = 0
text = '''
hostname --fqdn 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not Found'
fi
'''

[[metrics]]

id = "uname"
level = 0
text = '''
uname -a
'''

[[metrics]]

id = "lsb"
level = 0
text = '''
cat /etc/*-release
'''

[[metrics]]

id = "fips"
level = 1
text = '''
cat /proc/sys/crypto/fips_enabled
'''

[[metrics]]

id = "virturization"
level = 1
text = '''
cat /proc/cpuinfo
'''

[[metrics]]

id = "sestatus"
level = 1
text = '''
/usr/sbin/sestatus
'''

[[metrics]]

id = "mount_iso"
level = 1
text = '''
mount
'''

[[metrics]]

id = "proxy_global"
level = 1
text = '''
grep proxy /etc/yum.conf
if [ \$? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "kdump"
level = 1
text = '''
if [ -f /usr/bin/systemctl ]; then
    /usr/bin/systemctl status kdump
else
    /sbin/chkconfig --list|grep kdump
fi
'''

[[metrics]]

id = "crash_size"
level = 1
text = '''
cat /sys/kernel/kexec_crash_size 2>/dev/null
if [ $? != 0 ]; then
    echo 'Unkown crash_size. kdump:'
    cat /sys/kernel/kexec_crash_loaded
fi
'''

[[metrics]]

id = "kdump_path"
level = 1
text = '''
egrep -e '^(path|core_collector)' /etc/kdump.conf 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "iptables"
level = 1
text = '''
if [ -f /usr/bin/systemctl ]; then
    /usr/bin/systemctl status iptables firewalld
else
    /sbin/chkconfig --list|grep iptables
fi
'''

[[metrics]]

id = "runlevel"
level = 1
text = '''
if [ -f /usr/bin/systemctl ]; then
    /usr/bin/systemctl get-default
else
    grep :initdefault /etc/inittab
fi
'''

[[metrics]]

id = "resolve_conf"
level = 1
text = '''
grep nameserver /etc/resolv.conf 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not Found'
fi
'''

[[metrics]]

id = "keyboard"
level = 1
text = '''
if [ -f /etc/sysconfig/keyboard ]; then
    cat /etc/sysconfig/keyboard
elif [ -f /etc/vconsole.conf ]; then
    cat /etc/vconsole.conf
fi
'''

[[metrics]]

id = "language"
level = 1
text = '''
cat /proc/cmdline
'''

[[metrics]]

id = "grub"
level = 1
text = '''
grep GRUB_CMDLINE_LINUX /etc/default/grub 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "timezone"
level = 1
text = '''
if [ -x /bin/timedatectl ]; then
    /bin/timedatectl
elif [ -f /etc/sysconfig/clock ]; then
    cat /etc/sysconfig/clock
fi
'''

[[metrics]]

id = "ntp_slew"
level = 1
text = '''
grep -i options /etc/sysconfig/ntpd 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "ntp"
level = 1
text = '''
egrep -e '^server' /etc/ntp.conf 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "snmp_trap"
level = 1
text = '''
cat /etc/snmp/snmpd.conf
'''

[[metrics]]

id = "vmware_scsi_timeout"
level = 1
text = '''
cat /etc/udev/rules.d/99-vmware-scsi-udev.rules 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "vmwaretool_timesync"
level = 1
text = '''
LANG=c /usr/bin/vmware-toolbox-cmd timesync status 2>/dev/null
if [ $? == 127 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "cpu"
level = 0
text = '''
cat /proc/cpuinfo
'''

[[metrics]]

id = "meminfo"
level = 0
text = '''
cat /proc/meminfo
'''

[[metrics]]

id = "net_onboot"
level = 0
text = '''
cd /etc/sysconfig/network-scripts/
grep ONBOOT ifcfg-*
'''

[[metrics]]

id = "net_route"
level = 0
text = '''
/sbin/ip route
'''

[[metrics]]

id = "net_bond"
level = 0
text = '''
cd /etc/sysconfig/network-scripts/
cat *-bond* 2>/dev/null
if [ $? != 0 ]; then
    echo 'Not found'
fi
'''

[[metrics]]

id = "network"
level = 0
text = '''
/sbin/ip addr
'''

[[metrics]]

id = "block_device"
level = 0
text = '''
egrep -H '.*' /sys/block/*/size
egrep -H '.*' /sys/block/*/removable
egrep -H '.*' /sys/block/*/device/model
egrep -H '.*' /sys/block/*/device/rev
egrep -H '.*' /sys/block/*/device/state
egrep -H '.*' /sys/block/*/device/timeout
egrep -H '.*' /sys/block/*/device/vendor
egrep -H '.*' /sys/block/*/device/queue_depth
'''

[[metrics]]

id = "mdadb"
level = 0
text = '''
cat /proc/mdstat
'''

[[metrics]]

id = "fstab"
level = 0
text = '''
cat /etc/fstab
'''

[[metrics]]

id = "lvm"
level = 0
text = '''
mount
'''

[[metrics]]

id = "filesystem"
level = 0
text = '''
if [ -x /bin/lsblk ]; then
    /bin/lsblk -i
else
    /bin/df -k
fi
'''

[[metrics]]

id = "filesystem_df_ip"
level = 0
text = '''
df -iP
'''

[[metrics]]

id = "user"
level = 1
text = '''
cat /etc/passwd
'''

[[metrics]]

id = "group"
level = 1
deviceFlag = true
text = '''
cat /etc/group
'''

[[metrics]]

id = "service"
level = 1
text = '''
if [ -f /usr/bin/systemctl ]; then
    /usr/bin/systemctl list-units --type service --all
elif [ -f /sbin/chkconfig ]; then
    /sbin/chkconfig --list
fi
'''

[[metrics]]

id = "packages"
level = 1
text = '''
rpm -qa --qf "%{NAME}\t%|EPOCH?{%{EPOCH}}:{0}|\t%{VERSION}\t%{RELEASE}\t%{INSTALLTIME}\t%{ARCH}\n"
'''

[[metrics]]

id = "cron"
level = 1
text = '''
sudo -A sh -c "cd /var/spool/cron/; egrep -H '.*' *"
'''

[[metrics]]

id = "yum"
level = 1
text = '''
egrep -e '\[|enabled' /etc/yum.repos.d/*.repo
'''

[[metrics]]

id = "resource_limits"
level = 1
text = '''
egrep -v '^#' /etc/security/limits.d/*
'''

[[metrics]]

id = "error_messages"
level = 2
text = '''
egrep -i '(error|warning|failed)' /var/log/messages | head -100
'''

[[metrics]]

id = "oracle_module"
level = 2
text = '''
sudo -A ls /root/package/*
'''

[[metrics]]

id = "oracle"
level = 2
text = '''
ls -d /opt/oracle/app/product/*/* /*/app/oracle/product/*/* 2>/dev/null
if [ \$? != 0 ]; then
    echo 'Not found'
fi
'''
`

func (e *Linux) Label() string {
	return "Linux"
}

func (e *Linux) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("linuxconf", func() Exporter {
		host, _ := GetHostname()
		return &Linux{
			Server: host,
		}
	})
}
