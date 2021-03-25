package hpiloconf

type Metric struct {
	Level int    `toml:"level"`
	Id    string `toml:"id"`
	Text  string `toml:"text"`
}

type Metrics struct {
	Metrics []*Metric
}

func NewMetric(level int, id string, text string) *Metric {
	metric := &Metric{
		Level: level,
		Id:    id,
		Text:  text,
	}
	return metric
}

var metrics = []*Metric{
	NewMetric(0, "overview", "/json/overview"),
	NewMetric(0, "license", "/redfish/v1/Managers/1/LicenseService/1/"),
	NewMetric(0, "proc_info", "/json/proc_info"),
	NewMetric(0, "mem_info", "/json/mem_info"),
	NewMetric(0, "network", "/redfish/v1/Managers/1/EthernetInterfaces/1"),
	NewMetric(0, "health_drives", "/json/health_drives"),
	NewMetric(0, "health_phy_drives", "/json/health_phy_drives"),
	NewMetric(0, "snmp", "/redfish/v1/Managers/1/snmpservice/snmpalertdestinations/1"),
	NewMetric(0, "power_regulator", "/json/power_regulator"),
	NewMetric(0, "power_summary", "/json/power_summary"),
}
