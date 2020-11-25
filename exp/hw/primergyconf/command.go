package primergyconf

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
	NewCommand(0, "overview", "/json/overview"),
	// NewCommand(0, "license", "/redfish/v1/Managers/1/LicenseService/1/"),
	// NewCommand(0, "proc_info", "/json/proc_info"),
	// NewCommand(0, "mem_info", "/json/mem_info"),
	// NewCommand(0, "network", "/redfish/v1/Managers/1/EthernetInterfaces/1"),
	// NewCommand(0, "health_drives", "/json/health_drives"),
	// NewCommand(0, "health_phy_drives", "/json/health_phy_drives"),
	// NewCommand(0, "snmp", "/redfish/v1/Managers/1/snmpservice/snmpalertdestinations/1"),
	// NewCommand(0, "power_regulator", "/json/power_regulator"),
	// NewCommand(0, "power_summary", "/json/power_summary"),
}
