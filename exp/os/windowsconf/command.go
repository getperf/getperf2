package windowsconf

// type Env struct {
// 	Level      int
// 	DryRun     bool
// 	Timeout    int
// 	Datastore  string
// 	ConfigPath string
// 	LocalExec  bool
// 	Messages   string
// }

type Command struct {
	Level int    `toml:"level"`
	Id    string `toml:"id"`
	Text  string `toml:"text"`
}

type Commands struct {
	Commands []*Command
}

func NewCommand(level int, id string, text string) *Command {
	command := &Command{
		Level: level,
		Id:    id,
		Text:  text,
	}
	return command
}

var commands = []*Command{
	NewCommand(0, "system", `
			Get-WmiObject -Class Win32_ComputerSystem
			`),
	NewCommand(0, "os_conf", `
			Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion" | 
			Format-List
			`),
	NewCommand(0, "system", `
			Get-WmiObject -Class Win32_ComputerSystem
			`),
	NewCommand(0, "os_conf", `
			Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion" |
			Format-List
			`),
	NewCommand(0, "os", `
		Get-WmiObject Win32_OperatingSystem |
		Format-List Caption,CSDVersion,ProductType,OSArchitecture
		`),
	NewCommand(1, "driver", `
		Get-WmiObject Win32_PnPSignedDriver
		`),
	NewCommand(1, "fips", `
		Get-Item "HKLM:System\CurrentControlSet\Control\Lsa\FIPSAlgorithmPolicy"
		`),
	NewCommand(1, "virturalization", `
		Get-WmiObject -Class Win32_ComputerSystem | Select Model | FL
		`),
	NewCommand(1, "storage_timeout", `
		Get-ItemProperty "HKLM:SYSTEM\CurrentControlSet\Services\disk"
		`),
	NewCommand(1, "monitor", `
		Get-WmiObject Win32_DesktopMonitor | FL
		`),
	NewCommand(1, "ie_version", `
		Get-ItemProperty "HKLM:SOFTWARE\Microsoft\Internet Explorer"
		`),
	NewCommand(2, "system_log", `
		Get-EventLog system | Where-Object { $_.EntryType -eq "Error" } | FL
		`),
	NewCommand(2, "apps_log", `
		Get-EventLog application | Where-Object { $_.EntryType -eq "Error" } | FL
		`),
	NewCommand(1, "ntp", `
		(Get-Item "HKLM:System\CurrentControlSet\Services\W32Time\Parameters").GetValue("NtpServer")
		`),
	NewCommand(1, "user_account_control", `
		Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System"
		`),
	NewCommand(1, "remote_desktop", `
		(Get-Item "HKLM:System\CurrentControlSet\Control\Terminal Server").GetValue("fDenyTSConnections")
		`),
	NewCommand(0, "cpu", `
		Get-WmiObject -Class Win32_Processor | Format-List DeviceID, Name, MaxClockSpeed, SocketDesignation, NumberOfCores, NumberOfLogicalProcessors
		`),
	NewCommand(0, "memory", `
		Get-WmiObject Win32_OperatingSystem |
		select TotalVirtualMemorySize,TotalVisibleMemorySize,
		FreePhysicalMemory,FreeVirtualMemory,FreeSpaceInPagingFiles
		`),
	NewCommand(1, "dns", `
		Get-DnsClientServerAddress|FL
		`),
	NewCommand(2, "nic_teaming_config", `
		Get-NetLbfoTeamNic
		`),
	NewCommand(2, "tcp", `
		Get-ItemProperty "HKLM:SYSTEM\CurrentControlSet\Services\Tcpip\Parameters" |
		Format-List
		`),
	NewCommand(2, "network", `
		Get-WmiObject Win32_NetworkAdapterConfiguration |
		Where{$_.IpEnabled -Match "True"} |
		Select ServiceName, MacAddress, IPAddress, DefaultIPGateway, Description, IPSubnet |
		Format-List
		`),
	NewCommand(2, "network_profile", `
		Get-NetConnectionProfile
		`),
	NewCommand(2, "net_bind", `
		Get-NetAdapterBinding | FL
		`),
	NewCommand(2, "net_ip", `
		Get-NetIPInterface | FL
		`),
	NewCommand(2, "firewall", `
		Get-NetFirewallRule -Direction Inbound -Enabled True
		`),
	NewCommand(1, "filesystem", `
		Get-WmiObject Win32_LogicalDisk | Format-List *
		`),
	NewCommand(2, "user", `
		$result = @()
		$accountObjList =  Get-CimInstance -ClassName Win32_Account
		$userObjList = Get-CimInstance -ClassName Win32_UserAccount
		foreach($userObj in $userObjList)
		{
		    $IsLocalAccount = ($userObjList | ?{$_.SID -eq $userObj.SID}).LocalAccount
		    if($IsLocalAccount)
		    {
		        $query = "WinNT://{0}/{1},user" -F $env:COMPUTERNAME,$userObj.Name
		        $dirObj = New-Object -TypeName System.DirectoryServices.DirectoryEntry -ArgumentList $query
		        $UserFlags = $dirObj.InvokeGet("UserFlags")
		        $DontExpirePasswd = [boolean]($UserFlags -band 0x10000)
		        $AccountDisable   = [boolean]($UserFlags -band 0x2)
		        $obj = New-Object -TypeName PsObject
		        Add-Member -InputObject $obj -MemberType NoteProperty -Name "UserName" -Value $userObj.Name
		        Add-Member -InputObject $obj -MemberType NoteProperty -Name "DontExpirePasswd" -Value $DontExpirePasswd
		        Add-Member -InputObject $obj -MemberType NoteProperty -Name "AccountDisable" -Value $AccountDisable
		        Add-Member -InputObject $obj -MemberType NoteProperty -Name "SID" -Value $userObj.SID
		        $result += $obj
		    }
		}
		$result | Format-List
		`),
	NewCommand(1, "etc_hosts", `
		Get-Content "$($env:windir)\system32\Drivers\etc\hosts"
		`),
	NewCommand(0, "patch_lists", `
		wmic qfe
		`),
	NewCommand(1, "service", `
		Get-Service | FL
		`),
	NewCommand(0, "packages", `
		Get-WmiObject Win32_Product |
		Select-Object Name, Vendor, Version |
		Format-List
		Get-ChildItem -Path(
		'HKLM:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall',
		'HKCU:SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall') |
		% { Get-ItemProperty $_.PsPath | Select-Object DisplayName, Publisher, DisplayVersion } |
		Format-List
		`),
	NewCommand(1, "feature", `
		Get-WindowsFeature | ?{$_.InstallState -eq [Microsoft.Windows.ServerManager.Commands.InstallState]::Installed} | FL
		`),
	NewCommand(1, "task_scheduler", `
		Get-ScheduledTask |
		? {$_.State -eq "Ready"} |
		Get-ScheduledTaskInfo |
		? {$_.NextRunTime -ne $null}|
		Format-List
		`),
}
