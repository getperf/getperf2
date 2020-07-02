package linuxconf

import "io"

type ExecType string

const (
	Cmd    = ExecType("Cmd")
	Script = ExecType("Script")
)

type Command struct {
	Level int      `toml:"level"`
	Type  ExecType `toml:"type"`
	Id    string   `toml:"id"`
	Text  string   `toml:"text"`

	stdOut io.Writer
	stdErr io.Writer
}

type Commands struct {
	Commands []*Command
}

func NewCommand(level int, execType ExecType, id string, text string) *Command {
	command := &Command{
		Level: level,
		Type:  execType,
		Id:    id,
		Text:  text,
	}
	return command
}

var commands = []*Command{
	NewCommand(0, "Cmd",
		"hostname",
		`hostname -s`),
	NewCommand(0, "Script",
		"hostname_fqdn",
		`hostname --fqdn 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not Found'
		 fi`),
	NewCommand(0, "Cmd",
		"uname",
		`uname -a`),
	NewCommand(0, "Cmd",
		"lsb",
		`cat /etc/*-release`),
	NewCommand(0, "Cmd", "cpu",
		`cat /proc/cpuinfo`),
	NewCommand(0, "Script",
		"machineid",
		`if [ -f /etc/machine-id ]; then
		     cat /etc/machine-id
		 elif [ -f /var/lib/dbus/machine-id ]; then
		     cat /var/lib/dbus/machine-id
		 fi`),
	NewCommand(0, "Cmd",
		"meminfo",
		`cat /proc/meminfo`),
	NewCommand(0, "Cmd",
		"network",
		`/sbin/ip addr`),
	NewCommand(0, "Script",
		"net_onboot",
		`cd /etc/sysconfig/network-scripts/
		 grep ONBOOT ifcfg-*`),
	NewCommand(0, "Cmd",
		"net_route",
		`/sbin/ip route`),
	NewCommand(0, "Script",
		"net_bond",
		`cd /etc/sysconfig/network-scripts/
		 cat *-bond* 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"block_device",
		`egrep -H '.*' /sys/block/*/size
		 egrep -H '.*' /sys/block/*/removable
		 egrep -H '.*' /sys/block/*/device/model
		 egrep -H '.*' /sys/block/*/device/rev
		 egrep -H '.*' /sys/block/*/device/state
		 egrep -H '.*' /sys/block/*/device/timeout
		 egrep -H '.*' /sys/block/*/device/vendor
		 egrep -H '.*' /sys/block/*/device/queue_depth`),
	NewCommand(0, "Cmd",
		"mdadb",
		`cat /proc/mdstat`),
	NewCommand(0, "Cmd",
		"filesystem",
		`cat /etc/fstab`),
	NewCommand(0, "Cmd",
		"fips",
		`cat /proc/sys/crypto/fips_enabled`),
	NewCommand(0, "Cmd",
		"virturization",
		`cat /proc/cpuinfo`),
	NewCommand(0, "Cmd",
		"packages",
		`rpm -qa --qf "%{NAME}\t%|EPOCH?{%{EPOCH}}:{0}|\t%{VERSION}\t%{RELEASE}\t%{INSTALLTIME}\t%{ARCH}\n"`),
	NewCommand(0, "Script",
		"cron",
		`sudo -A ls /var/spool/cron/ |cat`),
	NewCommand(0, "Cmd",
		"yum",
		`egrep -e '\[|enabled' /etc/yum.repos.d/*.repo`),
	NewCommand(0, "Cmd",
		"resource_limits",
		`egrep -v '^#' /etc/security/limits.d/*`),
	NewCommand(0, "Cmd",
		"user",
		`cat /etc/passwd`),
	NewCommand(0, "Cmd",
		"group",
		`cat /etc/group`),
	NewCommand(0, "Cmd",
		"crontab",
		`crontab -l`),
	NewCommand(0, "Script",
		"service",
		`if [ -f /usr/bin/systemctl ]; then
		     /usr/bin/systemctl list-units --type service --all
		 elif [ -f /sbin/chkconfig ]; then
		     /sbin/chkconfig --list
		 fi`),
	NewCommand(0, "Cmd",
		"mount_iso",
		`mount`),
	NewCommand(0, "Script",
		"oracle",
		`ls -d /opt/oracle/app/product/*/* /*/app/oracle/product/*/* 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"proxy_global",
		`grep proxy /etc/yum.conf
		 if [ \$? != 0 ]; then
		     echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"kdump",
		`if [ -f /usr/bin/systemctl ]; then
		     /usr/bin/systemctl status kdump
		 elif [ -f /sbin/chkconfig ]; then
		     /sbin/chkconfig --list|grep kdump
		 fi`),
	NewCommand(0, "Script",
		"crash_size",
		`cat /sys/kernel/kexec_crash_size 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Unkown crash_size. kdump:'
		    cat /sys/kernel/kexec_crash_loaded
		 fi`),
	NewCommand(0, "Script",
		"kdump_path",
		`egrep -e '^(path|core_collector)' /etc/kdump.conf 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"iptables",
		`if [ -f /usr/bin/systemctl ]; then
		     /usr/bin/systemctl status iptables firewalld
		 elif [ -f /sbin/chkconfig ]; then
		     /sbin/chkconfig --list|grep iptables
		 fi`),
	NewCommand(0, "Script",
		"runlevel",
		`if [ -f /usr/bin/systemctl ]; then
		     /usr/bin/systemctl get-default
		 else
		     grep :initdefault /etc/inittab
		 fi`),
	NewCommand(0, "Script",
		"resolve_conf",
		`grep nameserver /etc/resolv.conf 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not Found'
		 fi`),
	NewCommand(0, "Script",
		"grub",
		`grep GRUB_CMDLINE_LINUX /etc/default/grub 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"ntp",
		`egrep -e '^server' /etc/ntp.conf 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"ntp_slew",
		`grep -i options /etc/sysconfig/ntpd 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Cmd",
		"snmp_trap",
		`cat /etc/snmp/snmpd.conf`),
	NewCommand(0, "Script",
		"sestatus",
		`/usr/sbin/sestatus`),
	NewCommand(0, "Script",
		"keyboard",
		`if [ -f /etc/sysconfig/keyboard ]; then
		     cat /etc/sysconfig/keyboard
		 elif [ -f /etc/vconsole.conf ]; then
		     cat /etc/vconsole.conf
		 fi`),
	NewCommand(0, "Script",
		"vmwaretool_timesync",
		`LANG=c /usr/bin/vmware-toolbox-cmd timesync status 2>/dev/null
		 if [ $? == 127 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Script",
		"vmware_scsi_timeout",
		`cat /etc/udev/rules.d/99-vmware-scsi-udev.rules 2>/dev/null
		 if [ $? != 0 ]; then
		    echo 'Not found'
		 fi`),
	NewCommand(0, "Cmd",
		"language",
		`cat /proc/cmdline`),
	NewCommand(0, "Script",
		"timezone",
		`if [ -x /bin/timedatectl ]; then
		     /bin/timedatectl
		 elif [ -f /etc/sysconfig/clock ]; then
		     cat /etc/sysconfig/clock
		 fi`),
	NewCommand(1, "Script",
		"error_messages",
		`egrep -i '(error|warning|failed)' /var/log/messages | head -100`),
	NewCommand(0, "Cmd",
		"oracle_module",
		`sudo -A ls /root/package/*`),
	NewCommand(0, "Script",
		"vncserver",
		`if [ -f /usr/bin/systemctl ]; then
		     /usr/bin/systemctl status vncserver
		 else
		     /sbin/chkconfig --list|grep vncserver
		 fi`),
}
