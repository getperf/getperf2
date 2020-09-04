package windowsconf

import (
	"io"

	. "github.com/getperf/getperf2/common"
	. "github.com/getperf/getperf2/exp"
)

const Version = "0.1.4"

type Windows struct {
	Server    string    `toml:"server"`
	LocalExec bool      `toml:"local_exec"`
	Servers   []*Server `toml:"servers"`
	Metrics   []*Metric `toml:"metrics"`

	errFile    io.Writer
	datastore  string
	ScriptPath string
}

type Server struct {
	Server   string `toml:"server"`
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Insecure bool   `toml:"insecure"`
}

var sampleTemplateConfig = `
# Windows inventory exporter settings
# When collecting the inventory of Windows platform, execute it locally.
# Therefore, no template setting is required

local_exec = true

## Enter the information for login account of remote linux server
# 
## example:
# 
# [[servers]]
# 
# server = "centos8"    # server name
# url = "192.168.10.1"  # server address, example: 192.168.0.1 , 192.168.0.1:22
# user = "test_user"
# password = "P@ssword"
# ssh_key = ""          # ssh private key path, ignore if it not set

{{if ne .Url "" }}
[[servers]]

server = "{{ .Server }}"
url = "{{ .Url }}"
user = "{{ .User }}"
password = "{{ .Password }}"
insecure = true
{{end}}

## Describe the additional command list. Added to the default command list for
## Windows Inventory scenarios. The text parameter using escape codes such as
## '\"', '\\', See these example,
#
## example:
# 
# [[metrics]]
#
# id = "echo"    # unique key
# type = "Cmd"   # Cmd : cmd.exe -c "..." , Cmdlet : PowerShell -Command {...}
# level = 0      # command level [0-2]
# text = "echo 1"

[[metrics]]

id = "system"
level = 0
type = "Cmdlet"
text = '''
Get-WmiObject -Class Win32_ComputerSystem
'''

[[metrics]]

id = "os_conf"
level = 0
type = "Cmdlet"
text = '''
Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' |
Format-List
'''

[[metrics]]

id = "os"
level = 0
type = "Cmdlet"
text = '''
Get-WmiObject Win32_OperatingSystem |
Format-List Caption,CSDVersion,ProductType,OSArchitecture
'''

[[metrics]]

id = "driver"
level = 1
type = "Cmdlet"
text = '''
Get-WmiObject Win32_PnPSignedDriver
'''

[[metrics]]

id = "fips"
level = 1
type = "Cmdlet"
text = '''
Get-Item 'HKLM:System\CurrentControlSet\Control\Lsa\FIPSAlgorithmPolicy'
'''

[[metrics]]

id = "virturalization"
level = 1
type = "Cmdlet"
text = '''
Get-WmiObject -Class Win32_ComputerSystem | Select Model | FL
'''

[[metrics]]

id = "storage_timeout"
level = 1
type = "Cmdlet"
text = '''
Get-ItemProperty 'HKLM:SYSTEM\CurrentControlSet\Services\disk'
'''

[[metrics]]

id = "monitor"
level = 1
type = "Cmdlet"
text = '''
Get-WmiObject Win32_DesktopMonitor | FL
'''

[[metrics]]

id = "ie_version"
level = 1
type = "Cmdlet"
text = '''
Get-ItemProperty 'HKLM:SOFTWARE\Microsoft\Internet Explorer'
'''

[[metrics]]

id = "system_log"
level = 2
type = "Cmdlet"
text = '''
Get-EventLog system | Where-Object { $_.EntryType -eq 'Error' } | FL
'''

[[metrics]]

id = "apps_log"
level = 2
type = "Cmdlet"
text = '''
Get-EventLog application | Where-Object { $_.EntryType -eq 'Error' } | FL
'''

[[metrics]]

id = "ntp"
level = 1
type = "Cmdlet"
text = '''
(Get-Item 'HKLM:System\CurrentControlSet\Services\W32Time\Parameters').GetValue('NtpServer')
'''

[[metrics]]

id = "user_account_control"
level = 1
type = "Cmdlet"
text = '''
Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System'
'''

[[metrics]]

id = "remote_desktop"
level = 1
type = "Cmdlet"
text = '''
(Get-Item 'HKLM:System\CurrentControlSet\Control\Terminal Server').GetValue('fDenyTSConnections')
'''

[[metrics]]

id = "cpu"
level = 0
type = "Cmdlet"
text = '''
Get-WmiObject -Class Win32_Processor | Format-List DeviceID, Name, MaxClockSpeed, SocketDesignation, NumberOfCores, NumberOfLogicalProcessors
'''

[[metrics]]

id = "memory"
level = 0
type = "Cmdlet"
text = '''
Get-WmiObject Win32_OperatingSystem |
select TotalVirtualMemorySize,TotalVisibleMemorySize,
FreePhysicalMemory,FreeVirtualMemory,FreeSpaceInPagingFiles
'''

[[metrics]]

id = "dns"
level = 1
type = "Cmdlet"
text = '''
Get-DnsClientServerAddress|FL
'''

[[metrics]]

id = "nic_teaming_config"
level = 2
type = "Cmdlet"
text = '''
Get-NetLbfoTeamNic
'''

[[metrics]]

id = "tcp"
level = 2
type = "Cmdlet"
text = '''
Get-ItemProperty 'HKLM:SYSTEM\CurrentControlSet\Services\Tcpip\Parameters' |
Format-List
'''

[[metrics]]

id = "network"
level = 2
type = "Cmdlet"
text = '''
Get-WmiObject Win32_NetworkAdapterConfiguration |
Where{$_.IpEnabled -Match 'True'} |
Select ServiceName, MacAddress, IPAddress, DefaultIPGateway, Description, IPSubnet |
Format-List
'''

[[metrics]]

id = "network_profile"
level = 2
type = "Cmdlet"
text = '''
Get-NetConnectionProfile | FL
'''

[[metrics]]

id = "net_bind"
level = 2
type = "Cmdlet"
text = '''
Get-NetAdapterBinding | FL
'''

[[metrics]]

id = "net_ip"
level = 2
type = "Cmdlet"
text = '''
Get-NetIPInterface | FL
'''

[[metrics]]

id = "firewall"
level = 2
type = "Cmdlet"
text = '''
Get-NetFirewallRule -Direction Inbound -Enabled True
'''

[[metrics]]

id = "filesystem"
level = 1
type = "Cmdlet"
text = '''
Get-WmiObject Win32_LogicalDisk | Format-List *
'''

[[metrics]]

id = "user"
level = 2
type = "Cmdlet"
text = '''
$result = @()
$accountObjList =  Get-CimInstance -ClassName Win32_Account
$userObjList = Get-CimInstance -ClassName Win32_UserAccount
foreach($userObj in $userObjList)
{
    $IsLocalAccount = ($userObjList | ?{$_.SID -eq $userObj.SID}).LocalAccount
    if($IsLocalAccount)
    {
        $query = 'WinNT://{0}/{1},user' -F $env:COMPUTERNAME,$userObj.Name
        $dirObj = New-Object -TypeName System.DirectoryServices.DirectoryEntry -ArgumentList $query
        $UserFlags = $dirObj.InvokeGet('UserFlags')
        $DontExpirePasswd = [boolean]($UserFlags -band 0x10000)
        $AccountDisable   = [boolean]($UserFlags -band 0x2)
        $obj = New-Object -TypeName PsObject
        Add-Member -InputObject $obj -MemberType NoteProperty -Name 'UserName' -Value $userObj.Name
        Add-Member -InputObject $obj -MemberType NoteProperty -Name 'DontExpirePasswd' -Value $DontExpirePasswd
        Add-Member -InputObject $obj -MemberType NoteProperty -Name 'AccountDisable' -Value $AccountDisable
        Add-Member -InputObject $obj -MemberType NoteProperty -Name 'SID' -Value $userObj.SID
        $result += $obj
    }
}
$result | Format-List
'''

[[metrics]]

id = "whoami"
level = 0
type = "Cmd"
text = '''
whoami /user
'''

[[metrics]]

id = "etc_hosts"
level = 1
type = "Cmdlet"
text = '''
Get-Content 'C:\Windows\system32\Drivers\etc\hosts'
'''

[[metrics]]

id = "patch_lists"
level = 0
type = "Cmd"
text = '''
wmic qfe
'''

[[metrics]]

id = "net_accounts"
level = 2
type = "Cmd"
text = '''
net accounts
'''

[[metrics]]

id = "service"
level = 1
type = "Cmdlet"
text = '''
Get-Service | FL
'''

[[metrics]]

id = "packages"
level = 0
type = "Cmdlet"
text = '''
Get-WmiObject Win32_Product |
Select-Object Name, Vendor, Version |
Format-List
Get-ChildItem -Path(
'HKLM:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall',
'HKCU:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall') |
% { Get-ItemProperty $_.PsPath | Select-Object DisplayName, Publisher, DisplayVersion } |
Format-List
'''

[[metrics]]

id = "feature"
level = 1
type = "Cmdlet"
text = '''
Get-WindowsFeature | ?{$_.InstallState -eq [Microsoft.Windows.ServerManager.Commands.InstallState]::Installed} | FL
'''

[[metrics]]

id = "task_scheduler"
level = 1
type = "Cmdlet"
text = '''
Get-ScheduledTask |
? {$_.State -eq 'Ready'} |
Get-ScheduledTaskInfo |
? {$_.NextRunTime -ne $null}|
Format-List
'''

`

func (e *Windows) Label() string {
	return "Windows"
}

func (e *Windows) Config() string {
	return sampleTemplateConfig
}

func init() {
	AddExporter("windowsconf", func() Exporter {
		host, _ := GetHostname()
		return &Windows{
			Server: host,
		}
	})
}
