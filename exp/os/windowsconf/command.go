package windowsconf

type ExecType string

const (
	Cmd    = ExecType("Cmd")
	Cmdlet = ExecType("Cmdlet")
)

type Command struct {
	Level int      `toml:"level"`
	Type  ExecType `toml:"type"`
	Id    string   `toml:"id"`
	Text  string   `toml:"text"`
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
	NewCommand(0, "Cmdlet", "system", `
			Get-WmiObject -Class Win32_ComputerSystem
			`),
	NewCommand(0, "Cmdlet", "os_conf", `
			Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' |
			Format-List
			`),
	NewCommand(0, "Cmdlet", "os", `
		Get-WmiObject Win32_OperatingSystem |
		Format-List Caption,CSDVersion,ProductType,OSArchitecture
		`),
	NewCommand(1, "Cmdlet", "driver", `
		Get-WmiObject Win32_PnPSignedDriver
		`),
	NewCommand(1, "Cmdlet", "fips", `
		Get-Item 'HKLM:System\CurrentControlSet\Control\Lsa\FIPSAlgorithmPolicy'
		`),
	NewCommand(1, "Cmdlet", "virturalization", `
		Get-WmiObject -Class Win32_ComputerSystem | Select Model | FL
		`),
	NewCommand(1, "Cmdlet", "storage_timeout", `
		Get-ItemProperty 'HKLM:SYSTEM\CurrentControlSet\Services\disk'
		`),
	NewCommand(1, "Cmdlet", "monitor", `
		Get-WmiObject Win32_DesktopMonitor | FL
		`),
	NewCommand(1, "Cmdlet", "ie_version", `
		Get-ItemProperty 'HKLM:SOFTWARE\Microsoft\Internet Explorer'
		`),
	NewCommand(2, "Cmdlet", "system_log", `
		Get-EventLog system | Where-Object { $_.EntryType -eq 'Error' } | FL
		`),
	NewCommand(2, "Cmdlet", "apps_log", `
		Get-EventLog application | Where-Object { $_.EntryType -eq 'Error' } | FL
		`),
	NewCommand(1, "Cmdlet", "ntp", `
		(Get-Item 'HKLM:System\CurrentControlSet\Services\W32Time\Parameters').GetValue('NtpServer')
		`),
	NewCommand(1, "Cmdlet", "user_account_control", `
		Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System'
		`),
	NewCommand(1, "Cmdlet", "remote_desktop", `
		(Get-Item 'HKLM:System\CurrentControlSet\Control\Terminal Server').GetValue('fDenyTSConnections')
		`),
	NewCommand(0, "Cmdlet", "cpu", `
		Get-WmiObject -Class Win32_Processor | Format-List DeviceID, Name, MaxClockSpeed, SocketDesignation, NumberOfCores, NumberOfLogicalProcessors
		`),
	NewCommand(0, "Cmdlet", "memory", `
		Get-WmiObject Win32_OperatingSystem |
		select TotalVirtualMemorySize,TotalVisibleMemorySize,
		FreePhysicalMemory,FreeVirtualMemory,FreeSpaceInPagingFiles
		`),
	NewCommand(1, "Cmdlet", "dns", `
		Get-DnsClientServerAddress|FL
		`),
	NewCommand(2, "Cmdlet", "nic_teaming_config", `
		Get-NetLbfoTeamNic
		`),
	NewCommand(2, "Cmdlet", "tcp", `
		Get-ItemProperty 'HKLM:SYSTEM\CurrentControlSet\Services\Tcpip\Parameters' |
		Format-List
		`),
	NewCommand(2, "Cmdlet", "network", `
		Get-WmiObject Win32_NetworkAdapterConfiguration |
		Where{$_.IpEnabled -Match 'True'} |
		Select ServiceName, MacAddress, IPAddress, DefaultIPGateway, Description, IPSubnet |
		Format-List
		`),
	NewCommand(2, "Cmdlet", "network_profile", `
		Get-NetConnectionProfile
		`),
	NewCommand(2, "Cmdlet", "net_bind", `
		Get-NetAdapterBinding | FL
		`),
	NewCommand(2, "Cmdlet", "net_ip", `
		Get-NetIPInterface | FL
		`),
	NewCommand(2, "Cmdlet", "firewall", `
		Get-NetFirewallRule -Direction Inbound -Enabled True
		`),
	NewCommand(1, "Cmdlet", "filesystem", `
		Get-WmiObject Win32_LogicalDisk | Format-List *
		`),
	NewCommand(2, "Cmdlet", "user", `
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
		`),
	NewCommand(1, "Cmdlet", "etc_hosts", `
		Get-Content '$($env:windir)\system32\Drivers\etc\hosts'
		`),
	NewCommand(0, "Cmd", "patch_lists", `
		wmic qfe
		`),
	NewCommand(1, "Cmdlet", "service", `
		Get-Service | FL
		`),
	NewCommand(0, "Cmdlet", "packages", `
		Get-WmiObject Win32_Product |
		Select-Object Name, Vendor, Version |
		Format-List
		`),
	NewCommand(0, "Cmdlet", "packages2", `
		Get-ChildItem -Path(
		'HKLM:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall',
		'HKCU:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall') |
		% { Get-ItemProperty $_.PsPath | Select-Object DisplayName, Publisher, DisplayVersion } |
		Format-List
		`),
	NewCommand(1, "Cmdlet", "feature", `
		Get-WindowsFeature | ?{$_.InstallState -eq [Microsoft.Windows.ServerManager.Commands.InstallState]::Installed} | FL
		`),
	NewCommand(1, "Cmdlet", "task_scheduler", `
		Get-ScheduledTask |
		? {$_.State -eq 'Ready'} |
		Get-ScheduledTaskInfo |
		? {$_.NextRunTime -ne $null}|
		Format-List
		`),
}
